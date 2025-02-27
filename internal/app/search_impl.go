package app

import (
	"cmp"
	"context"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type searchService struct {
	db            store.DBTX
	queries       *store.Queries
	tagsService   TagsService
	uploadService *UploadService
	userService   UserService
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

func int32RangeToInt4Range(r Int32Range) store.Int4Range {
	return store.Int4Range{
		Min: pgtype.Int4{Int32: r.Min.Int32, Valid: r.Min.Valid},
		Max: pgtype.Int4{Int32: r.Max.Int32, Valid: r.Max.Valid},
	}
}

// SearchBooks implements SearchService.
func (s *searchService) SearchBooks(ctx context.Context, req BookSearchQuery) (*BookSearchResult, error) {
	dbReq, err := constructBookSearchRequest(ctx, s.tagsService, req)
	if err != nil {
		return nil, err
	}

	start := time.Now()

	result, err := s.searchInternal(ctx, dbReq)
	if err != nil {
		return nil, err
	}

	result.TookUS = time.Since(start).Microseconds()
	return result, nil

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

func (s *searchService) searchInternal(ctx context.Context, dbReq store.BookSearchRequest) (*BookSearchResult, error) {
	books, err := store.SearchBooks(ctx, s.db, dbReq)
	if err != nil {
		return nil, err
	}

	totalBooks, err := store.CountBooks(ctx, s.db, dbReq, 1_000_000)
	if err != nil {
		return nil, err
	}

	result := new(BookSearchResult)
	result.Page = uint32(dbReq.Page)
	result.PageSize = uint32(dbReq.PageSize)
	result.TotalPages = uint32(math.Ceil(float64(totalBooks) / float64(dbReq.PageSize)))
	result.Books = make([]BookSearchItem, len(books))
	bookIds := make([]int64, len(books))

	tagsAgg := newTagsAggregator(s.tagsService)

	for i, book := range books {
		bookIds[i] = book.ID
		result.Books[i] = BookSearchItem{
			ID:   book.ID,
			Name: book.Name,
			Author: BookDetailsAuthorDto{
				ID:   uuidDbToDomain(book.AuthorUserID),
				Name: book.AuthorName,
			},
			CreatedAt:       book.CreatedAt.Time,
			AgeRating:       ageRatingFromDbValue(book.AgeRating),
			Words:           int(book.Words),
			Chapters:        int(book.Chapters),
			WordsPerChapter: getWordsPerChapter(int(book.Words), int(book.Chapters)),
			Summary:         book.Summary,
			Cover:           getBookCoverURL(s.uploadService, book.ID, book.HasCover),
			Tags:            arrInt64ToInt64String(book.TagIds),
		}

		tagsAgg.Add(book.ID, book.TagIds)
	}

	{
		tags, err := tagsAgg.Fetch(ctx)
		if err != nil {
			return nil, err
		}

		result.Tags = make([]DefinedTagDto, len(tags))
		i := 0
		for _, tag := range tags {
			result.Tags[i] = tag
			i += 1
		}
		slices.SortFunc(result.Tags, func(a, b DefinedTagDto) int {
			return cmp.Compare[string](a.Name, b.Name)
		})
	}

	// collections
	{
		bookCollectionsArr, err := s.queries.GetBooksCollections(ctx, bookIds)
		if err != nil {
			return nil, err
		}
		bookCollectionMap := make(map[int64][]BookCollectionDto, 0)

		for _, bookCollection := range bookCollectionsArr {
			item := BookCollectionDto{
				ID:       bookCollection.ID,
				Name:     bookCollection.Name,
				Position: int(bookCollection.Position),
				Size:     int(bookCollection.Size),
			}

			if _, ok := bookCollectionMap[bookCollection.BookID]; !ok {
				bookCollectionMap[bookCollection.BookID] = []BookCollectionDto{item}
			} else {
				bookCollectionMap[bookCollection.BookID] = append(bookCollectionMap[bookCollection.BookID], item)
			}
		}

		for i := 0; i < len(result.Books); i++ {
			bookCollections, ok := bookCollectionMap[result.Books[i].ID]
			if ok {
				result.Books[i].Collections = bookCollections
			} else {
				result.Books[i].Collections = []BookCollectionDto{}
			}
		}
	}

	return result, nil
}

func constructBookSearchRequest(ctx context.Context, tagsService TagsService, req BookSearchQuery) (dbReq store.BookSearchRequest, err error) {
	query := strings.Trim(req.Query, " \n\t")
	if len(query) >= 2 {
		dbReq.Query = req.Query
	}

	{
		var includeTags BookTags
		includeTags, err = tagsService.FindParentTagIds(ctx, req.IncludeTags)
		if err != nil {
			return
		}
		dbReq.IncludeParentTags = includeTags.ParentTagIds
	}

	{
		var excludeTags BookTags
		excludeTags, err = tagsService.FindParentTagIds(ctx, req.ExcludeTags)
		if err != nil {
			return
		}
		dbReq.ExcludeParentTags = excludeTags.ParentTagIds
	}

	dbReq.Words = int32RangeToInt4Range(req.Words)
	dbReq.WordsPerChapter = int32RangeToInt4Range(req.WordsPerChapter)
	dbReq.Chapters = int32RangeToInt4Range(req.Chapters)

	dbReq.IncludeAuthors = arrUuidDomainToDb(req.IncludeUsers)
	dbReq.ExcludeAuthors = arrUuidDomainToDb(req.ExcludeUsers)
	dbReq.IncludeBanned = req.IncludeBanned
	dbReq.IncludeHidden = req.IncludeHidden
	dbReq.IncludeEmpty = req.IncludeEmpty
	dbReq.Page = req.Page
	dbReq.PageSize = req.PageSize

	return
}

func NewSearchService(db store.DBTX, tagsService TagsService, uploadService *UploadService, userService UserService) SearchService {
	return &searchService{
		db:            db,
		queries:       store.New(db),
		tagsService:   tagsService,
		uploadService: uploadService,
		userService:   userService,
	}
}
