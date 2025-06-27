package app

import (
	"context"
	"math"
	"time"

	"github.com/MaratBR/openlibrary/internal/commonutil"
	elasticstore "github.com/MaratBR/openlibrary/internal/elastic-store"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/lib/gset"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/gofrs/uuid"
)

type searchService struct {
	db            store.DBTX
	queries       *store.Queries
	tagsService   TagsService
	uploadService *UploadService
	userService   UserService
	esClient      *elasticsearch.TypedClient
}

// ExplainSearchQuery implements SearchService.
func (s *searchService) ExplainSearchQuery(ctx context.Context, req BookSearchQuery) (DetailedBookSearchQuery, error) {
	detailed := DetailedBookSearchQuery{
		Query:           req.Query,
		Words:           req.Words,
		WordsPerChapter: req.WordsPerChapter,
		Chapters:        req.Chapters,
		IncludeBanned:   req.IncludeBanned,
		IncludeHidden:   req.IncludeHidden,
		IncludeEmpty:    req.IncludeEmpty,
		Page:            req.Page,
		PageSize:        req.PageSize,
	}

	tagIds := commonutil.MergeArrays(req.IncludeTags, req.ExcludeTags)

	if len(tagIds) == 0 {
		detailed.IncludeTags = []DefinedTagDto{}
		detailed.ExcludeTags = []DefinedTagDto{}
	} else {
		tags, err := s.tagsService.GetTagsByIds(ctx, tagIds)
		if err != nil {
			return DetailedBookSearchQuery{}, err
		}
		tagMap := make(map[int64]DefinedTagDto, len(tags))
		for _, t := range tags {
			tagMap[t.ID] = t
		}

		detailed.IncludeTags = []DefinedTagDto{}
		detailed.ExcludeTags = []DefinedTagDto{}

		for _, tagId := range req.IncludeTags {
			if t, ok := tagMap[tagId]; ok {
				detailed.IncludeTags = append(detailed.IncludeTags, t)
			}
		}

		for _, tagId := range req.ExcludeTags {
			if t, ok := tagMap[tagId]; ok {
				detailed.ExcludeTags = append(detailed.ExcludeTags, t)
			}
		}
	}

	userIds := commonutil.MergeArraysNoSort(req.IncludeUsers, req.ExcludeUsers)

	if len(userIds) == 0 {
		detailed.IncludeUsers = []UserFromSearchRequestDto{}
		detailed.ExcludeUsers = []UserFromSearchRequestDto{}
	} else {
		userMap := map[uuid.UUID]UserFromSearchRequestDto{}

		for _, userId := range userIds {
			user, err := s.userService.GetUserSelfData(ctx, userId)
			if err != nil {
				return DetailedBookSearchQuery{}, err
			}

			userMap[user.ID] = UserFromSearchRequestDto{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: getUserAvatar(user.Name, 84),
			}
		}

		detailed.IncludeUsers = []UserFromSearchRequestDto{}
		detailed.ExcludeUsers = []UserFromSearchRequestDto{}

		for _, userId := range req.IncludeUsers {
			if u, ok := userMap[userId]; ok {
				detailed.IncludeUsers = append(detailed.IncludeUsers, u)
			}
		}

		for _, userId := range req.ExcludeUsers {
			if u, ok := userMap[userId]; ok {
				detailed.ExcludeUsers = append(detailed.ExcludeUsers, u)
			}
		}
	}

	return detailed, nil
}

// SearchBooks implements SearchService.
func (s *searchService) SearchBooks(ctx context.Context, req BookSearchQuery) (*BookSearchResult, error) {
	// convert to elastic request
	esReq := elasticstore.SearchRequest{
		Query:        req.Query,
		IncludeUsers: make([]string, len(req.IncludeUsers)),
		ExcludeUsers: make([]string, len(req.ExcludeUsers)),
		IncludeTags:  req.IncludeTags,
		ExcludeTags:  req.ExcludeTags,
		Words: elasticstore.Range{
			Min: elasticstore.Int32{Int32: req.Words.Min.Int32, Valid: req.Words.Min.Valid},
			Max: elasticstore.Int32{Int32: req.Words.Max.Int32, Valid: req.Words.Max.Valid},
		},
		WordsPerChapter: elasticstore.Range{
			Min: elasticstore.Int32{Int32: req.WordsPerChapter.Min.Int32, Valid: req.WordsPerChapter.Min.Valid},
			Max: elasticstore.Int32{Int32: req.WordsPerChapter.Max.Int32, Valid: req.WordsPerChapter.Max.Valid},
		},
		Chapters: elasticstore.Range{
			Min: elasticstore.Int32{Int32: req.Chapters.Min.Int32, Valid: req.Chapters.Min.Valid},
			Max: elasticstore.Int32{Int32: req.Chapters.Max.Int32, Valid: req.Chapters.Max.Valid},
		},
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}

	now := time.Now()
	result, err := elasticstore.Search(ctx, s.esClient, esReq)
	tookTotal := time.Since(now).Microseconds()
	if err != nil {
		return nil, err
	}

	totalPages := uint32(math.Ceil(float64(result.Total) / float64(req.PageSize)))

	// construct search result
	searchResult := &BookSearchResult{
		TookUSTotal: tookTotal,
		TookUS:      result.TookMS * 1000,
		Books:       make([]BookSearchItem, len(result.Hits)),
		PageSize:    uint32(req.PageSize),
		Page:        uint32(req.Page),
		TotalPages:  totalPages,
		Total:       result.Total,
	}

	// get aggregated data

	var (
		bookIds     []int64
		authorNames map[uuid.UUID]string
		tags        map[int64]DefinedTagDto
		booksDbData map[int64]store.GetBookSearchRelatedDataRow
	)

	{
		booksDbData = make(map[int64]store.GetBookSearchRelatedDataRow, len(result.Hits))
		bookIds = make([]int64, len(result.Hits))
		authorIdsSet := gset.New[uuid.UUID]()
		tagsAgg := newTagsAggregator(s.tagsService)

		for i, book := range result.Hits {
			authorIdsSet.Add(book.AuthorID)
			bookIds[i] = book.ID
			tagsAgg.Add(book.ID, book.Tags)
		}

		queries := store.New(s.db)

		{
			authorIds := authorIdsSet.Arr()
			authorNamesRows, err := queries.GetUserNames(ctx, arrUuidDomainToDb(authorIds))
			if err != nil {
				return nil, wrapUnexpectedDBError(err)
			}
			authorNames = make(map[uuid.UUID]string, len(authorNamesRows))
			for _, row := range authorNamesRows {
				authorNames[uuidDbToDomain(row.ID)] = row.Name
			}
		}

		{
			rows, err := queries.GetBookSearchRelatedData(ctx, bookIds)
			if err != nil {
				return nil, wrapUnexpectedDBError(err)
			}
			for _, row := range rows {
				booksDbData[row.ID] = row
			}
		}

		tags, err = tagsAgg.Fetch(ctx)

		if err != nil {
			return nil, err
		}
	}

	for i, book := range result.Hits {
		authorName, _ := authorNames[book.AuthorID]
		bookData, _ := booksDbData[book.ID]

		searchResult.Books[i] = BookSearchItem{
			ID:   book.ID,
			Name: book.Name,
			Author: BookDetailsAuthorDto{
				ID:   book.AuthorID,
				Name: authorName,
			},
			CreatedAt:       bookData.CreatedAt.Time,
			AgeRating:       ageRatingFromDbValue(store.AgeRating(book.Rating)),
			Words:           int(book.Words),
			Chapters:        int(book.Chapters),
			WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
			Summary:         book.Description,
			Cover:           getBookCoverURL(s.uploadService, book.ID, bookData.HasCover),
			Tags:            arrInt64ToInt64String(book.Tags),
		}
	}

	// add tags
	searchResult.Tags = make([]DefinedTagDto, 0, len(tags))
	for _, tag := range tags {
		searchResult.Tags = append(searchResult.Tags, tag)
	}

	return searchResult, nil
}

func (s *searchService) GetBookExtremes(ctx context.Context) (*BookExtremes, error) {
	result, err := store.GetBooksFilterExtremes(ctx, s.db, &store.BookSearchFilter{
		IncludeBanned: false,
		IncludeHidden: false,
		IncludeEmpty:  false,
	})
	if err != nil {
		return nil, err
	}

	return &BookExtremes{
		Words:           Int32Range{Min: Int32{Int32: int32(result.WordsMin), Valid: true}, Max: Int32{Int32: int32(result.WordsMax), Valid: true}},
		Chapters:        Int32Range{Min: Int32{Int32: int32(result.ChaptersMin), Valid: true}, Max: Int32{Int32: int32(result.ChaptersMax), Valid: true}},
		WordsPerChapter: Int32Range{Min: Int32{Int32: int32(result.WordsPerChapterMin), Valid: true}, Max: Int32{Int32: int32(result.WordsPerChapterMax), Valid: true}},
		Favorites:       Int32Range{Min: Int32{Int32: int32(result.FavoritesMin), Valid: true}, Max: Int32{Int32: int32(result.FavoritesMax), Valid: true}},
	}, nil
}

func NewSearchService(db store.DBTX, tagsService TagsService, uploadService *UploadService, userService UserService, esClient *elasticsearch.TypedClient) SearchService {
	return &searchService{
		db:            db,
		queries:       store.New(db),
		tagsService:   tagsService,
		uploadService: uploadService,
		userService:   userService,
		esClient:      esClient,
	}
}
