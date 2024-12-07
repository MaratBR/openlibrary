package server

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/MaratBR/openlibrary/cmd/server/olproto"
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/commonutil"
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

type searchRequest struct {
	IncludeUsers    []uuid.UUID
	ExcludeUsers    []uuid.UUID
	IncludeTags     []int64
	ExcludeTags     []int64
	Words           app.Int32Range
	Chapters        app.Int32Range
	WordsPerChapter app.Int32Range
	Favorites       app.Int32Range
	FromBody        bool

	Page  int32
	Limit int32
}

type searchJSONResponse struct {
	Tags        []app.DefinedTagDto      `json:"tags"`
	BooksTookUS int64                    `json:"booksTook"`
	BooksMeta   app.BookSearchResultMeta `json:"booksMeta"`
	Books       []app.BookSearchItem     `json:"books"`
}

func (c *searchController) Search(w http.ResponseWriter, r *http.Request) {
	search, err := getSearchRequest(r)
	if err != nil {
		writeRequestError(err, w)
		return
	}

	pbResponse, err := performBookSearch(c.searchService, c.tagsService, r, search)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeProtobuf(w, pbResponse)
}

func getSearchRequest(r *http.Request) (search searchRequest, err error) {
	query := r.URL.Query()

	var source url.Values

	if query.Get("body") == "true" {
		source, err = readUrlEncodedBody(r)
		if err != nil {
			return
		}
		search.FromBody = true
	} else {
		source = query
	}

	search.Words = getInt32RangeFromQuery(source, "w")
	search.Favorites = getInt32RangeFromQuery(source, "f")
	search.Chapters = getInt32RangeFromQuery(source, "c")
	search.WordsPerChapter = getInt32RangeFromQuery(source, "wc")
	search.IncludeTags = getInt64Array(source, "it")
	search.ExcludeTags = getInt64Array(source, "et")
	search.IncludeUsers = getUUIDArray(source, "iu")
	search.ExcludeUsers = getUUIDArray(source, "eu")

	search.Page = getPage(source, "p")
	search.Limit = 20

	return
}

func performBookSearch(searchService app.SearchService, tagsService app.TagsService, r *http.Request, search searchRequest) (*olproto.ProtoSearchResult, error) {

	var (
		offset uint
		limit  uint
	)

	if search.Limit < 0 {
		limit = 20
	} else {
		limit = uint(search.Limit)
	}

	if search.Page > 0 {
		offset = uint(search.Page-1) * limit
	} else {
		offset = 0
	}

	result, err := searchService.SearchBooks(r.Context(), app.BookSearchQuery{
		UserID:          getNullableUserID(r),
		IncludeUsers:    search.IncludeUsers,
		ExcludeUsers:    search.ExcludeUsers,
		IncludeTags:     search.IncludeTags,
		ExcludeTags:     search.ExcludeTags,
		Words:           search.Words,
		Chapters:        search.Chapters,
		WordsPerChapter: search.WordsPerChapter,
		IncludeBanned:   false,
		IncludeHidden:   false,
		IncludeEmpty:    false,
		Offset:          offset,
		Limit:           limit,
	})

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

	// response := searchJSONResponse{
	// 	Books:       result.Books,
	// 	BooksMeta:   result.Meta,
	// 	BooksTookUS: result.TookUS,
	// 	Tags:        result.Tags,
	// }

	pbResponse := &olproto.ProtoSearchResult{
		Took:       uint32(result.TookUS),
		Page:       uint32(search.Page),
		CacheKey:   result.Meta.CacheKey,
		CacheTook:  uint32(result.Meta.CacheTookUS),
		CacheHit:   result.Meta.CacheHit,
		TotalPages: 0,
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
