package server

import (
	"net/http"
	"net/url"

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
	IncludeTags     []string
	ExcludeTags     []string
	Words           app.Int32Range
	Chapters        app.Int32Range
	WordsPerChapter app.Int32Range
	Favorites       app.Int32Range
	FromBody        bool
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
	search.IncludeTags = getStringArray(source, "it")
	search.ExcludeTags = getStringArray(source, "et")
	search.IncludeUsers = getUUIDArray(source, "iu")
	search.ExcludeUsers = getUUIDArray(source, "eu")

	return
}

func (c *searchController) performSearch(r *http.Request, search searchRequest) (*app.BookSearchResult, error) {
	return c.searchService.SearchBooks(r.Context(), app.BookSearchQuery{
		UserID: getNullableUserID(r),

		IncludeUsers: search.IncludeUsers,
		ExcludeUsers: search.ExcludeUsers,

		IncludeTags: search.IncludeTags,
		ExcludeTags: search.ExcludeTags,

		Words:           search.Words,
		Chapters:        search.Chapters,
		WordsPerChapter: search.WordsPerChapter,

		IncludeBanned: false,
		IncludeHidden: false,
		IncludeEmpty:  false,

		Offset: 0,
		Limit:  20,
	})
}

func (c *searchController) SearchPreload(r *http.Request, serverData *serverData) error {
	search, err := getSearchRequest(r)
	if err != nil {
		return err
	}
	result, err := c.performSearch(r, search)
	if err != nil {
		return err
	}

	{
		key := "/api/search?" + r.URL.RawQuery
		serverData.AddPreloadedData(key, result)

		for _, book := range result.Books {
			if book.Cover != "" {
				serverData.Preloads = append(serverData.Preloads, Preload{
					As:   "image",
					Href: book.Cover,
				})
			}
		}
		serverData.AddPreloadedData(key, result)
	}

	{
		bookExtremes, err := c.searchService.GetBookExtremes(r.Context())
		if err == nil {
			serverData.AddPreloadedData("/api/search/book-extremes", bookExtremes)
		}
	}

	if len(search.ExcludeTags)+len(search.IncludeTags) > 0 {
		tagNames := commonutil.MergeStringArrays(search.IncludeTags, search.ExcludeTags)
		tagIds, err := c.tagsService.FindBookTags(r.Context(), tagNames)
		if err == nil {
			tags, err := c.tagsService.GetTagsByIds(r.Context(), tagIds.TagIds)
			if err == nil {
				key := "/api/tags/lookup?q=" + url.QueryEscape(stringArray(tagNames))
				serverData.AddPreloadedData(key, tags)
			}
		}

	}

	return nil
}

func (c *searchController) Search(w http.ResponseWriter, r *http.Request) {
	search, err := getSearchRequest(r)
	if err != nil {
		writeRequestError(err, w)
		return
	}

	result, err := c.performSearch(r, search)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeJSON(w, result)
}

func (c *searchController) GetBookExtremes(w http.ResponseWriter, r *http.Request) {
	result, err := c.searchService.GetBookExtremes(r.Context())
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeJSON(w, result)
}
