package main

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/MaratBR/openlibrary/cmd/server/olproto"
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/gofrs/uuid"
)

type searchController struct {
	searchService app.SearchService
	tagsService   app.TagsService
}

func newSearchController(searchService app.SearchService, tagsService app.TagsService) *searchController {
	return &searchController{
		searchService: searchService,
		tagsService:   tagsService,
	}
}

func (c *searchController) Search(w http.ResponseWriter, r *http.Request) {
	search, err := getSearchRequest(r)
	if err != nil {
		writeBadRequest(err, w)
		return
	}

	pbResponse, err := performBookSearch(c.searchService, c.tagsService, r, search)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeProtobuf(w, pbResponse)
}

func getSearchRequest(r *http.Request) (searchRequest, error) {
	query := r.URL.Query()
	req := parseSearchRequest(query)
	return req, nil
}

func performBookSearch(searchService app.SearchService, tagsService app.TagsService, r *http.Request, search searchRequest) (*olproto.ProtoSearchResult, error) {

	result, err := searchService.SearchBooks(r.Context(), app.BookSearchQuery{
		UserID:          auth.GetNullableUserID(r.Context()),
		IncludeUsers:    search.IncludeUsers,
		ExcludeUsers:    search.ExcludeUsers,
		IncludeTags:     search.IncludeTags,
		ExcludeTags:     search.ExcludeTags,
		Words:           search.Words,
		Chapters:        search.Chapters,
		Favorites:       search.Favorites,
		WordsPerChapter: search.WordsPerChapter,

		IncludeBanned: false,
		IncludeHidden: false,
		IncludeEmpty:  false,
		Page:          search.Page,
		PageSize:      search.PageSize,
	})

	if err != nil {
		return nil, err
	}

	{
		allTagIds := commonutil.MergeArrays(search.IncludeTags, search.ExcludeTags)
		missingTagIds := []int64{}

		for _, tagId := range allTagIds {
			if !slices.ContainsFunc(result.Tags, func(tag app.DefinedTagDto) bool {
				return tag.ID == tagId
			}) {
				missingTagIds = append(missingTagIds, tagId)
			}
		}

		missingTags, err := tagsService.GetTagsByIds(r.Context(), allTagIds)
		if err != nil {
			return nil, err
		}

		result.Tags = append(result.Tags, missingTags...)
	}

	pbResponse := &olproto.ProtoSearchResult{
		Items: commonutil.MapSlice(result.Books, func(book app.BookSearchItem) *olproto.ProtoBookSearchItem {
			return &olproto.ProtoBookSearchItem{
				Id:         book.ID,
				Name:       book.Name,
				Cover:      book.Cover,
				AuthorId:   book.Author.ID.String(),
				AuthorName: book.Author.Name,
				AgeRating:  ageRatingToProto(book.AgeRating),
				Chapters:   uint32(book.Chapters),
				Favorites:  uint32(book.Favorites),
				Words:      uint32(book.Words),
				Summary:    book.Summary,
				TagIds:     commonutil.MapSlice(book.Tags, func(id app.Int64String) int64 { return int64(id) }),
			}
		}),
		Tags: commonutil.MapSlice(result.Tags, func(tag app.DefinedTagDto) *olproto.ProtoDefinedTag {
			return &olproto.ProtoDefinedTag{
				Id:          int64(tag.ID),
				Name:        tag.Name,
				Description: tag.Description,
				IsAdult:     tag.IsAdult,
				IsSpoiler:   tag.IsSpoiler,
				Category:    tagCategoryToProto(tag.Category),
			}
		}),
		TotalPages: result.TotalPages,
		PageSize:   result.PageSize,
		Page:       result.Page,
		Took:       uint32(result.TookUS),
		CacheKey:   result.Meta.CacheKey,
		CacheTook:  uint32(result.Meta.CacheTookUS),
		CacheHit:   result.Meta.CacheHit,
		Filter: &olproto.ProtoSearchFilter{
			WordsMin:           search.Words.Min.Ptr(),
			WordsMax:           search.Words.Max.Ptr(),
			ChaptersMin:        search.Chapters.Min.Ptr(),
			ChaptersMax:        search.Chapters.Max.Ptr(),
			WordsPerChapterMin: search.WordsPerChapter.Min.Ptr(),
			WordsPerChapterMax: search.WordsPerChapter.Max.Ptr(),
			FavoritesMin:       search.Favorites.Min.Ptr(),
			FavoritesMax:       search.Favorites.Max.Ptr(),
			IncludeTags:        search.IncludeTags,
			ExcludeTags:        search.ExcludeTags,
			IncludeUsers:       commonutil.StringifyUUIDArray(search.IncludeUsers),
			ExcludeUsers:       commonutil.StringifyUUIDArray(search.ExcludeUsers),
		},
	}

	return pbResponse, err
}

func (c *searchController) GetBookExtremes(w http.ResponseWriter, r *http.Request) {
	result, err := c.searchService.GetBookExtremes(r.Context())
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeJSON(w, result)
}

type searchRequest struct {
	IncludeUsers []uuid.UUID
	ExcludeUsers []uuid.UUID
	IncludeTags  []int64
	ExcludeTags  []int64

	Words           app.Int32Range
	Chapters        app.Int32Range
	WordsPerChapter app.Int32Range
	Favorites       app.Int32Range

	Page     uint
	PageSize uint
}

func parseSearchRequest(source url.Values) (search searchRequest) {
	search.Words = getInt32RangeFromQuery(source, "w")
	search.Favorites = getInt32RangeFromQuery(source, "f")
	search.Chapters = getInt32RangeFromQuery(source, "c")
	search.WordsPerChapter = getInt32RangeFromQuery(source, "wc")

	search.IncludeTags = getInt64Array(source, "it")
	search.ExcludeTags = getInt64Array(source, "et")

	search.IncludeUsers = getUUIDArray(source, "iu")
	search.ExcludeUsers = getUUIDArray(source, "eu")

	// pagination and page size
	search.Page = uint(olhttp.GetPage(source, "p"))
	pageSize := getInt32FromQuery(source, "ps")
	if pageSize.Valid {
		if pageSize.Int32 <= 0 {
			search.PageSize = 20
		} else if pageSize.Int32 > 100 {
			search.PageSize = 100
		} else {
			search.PageSize = uint(pageSize.Int32)
		}
	} else {
		search.PageSize = 20
	}

	return
}
