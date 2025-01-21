package publicui

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/web/public-ui/templates"
)

type searchController struct {
	service app.SearchService
}

func newSearchController(service app.SearchService) *searchController {
	return &searchController{
		service: service,
	}
}

func (c *searchController) Search(w http.ResponseWriter, r *http.Request) {
	search := getSearchRequest(r)
	query := app.BookSearchQuery{
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
	}
	result, err := c.service.SearchBooks(r.Context(), query)

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	explainedQuery, err := c.service.ExplainSearchQuery(r.Context(), query)

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	writeTemplate(w, r.Context(), templates.SearchPage(r.Context(), result, explainedQuery))
}
