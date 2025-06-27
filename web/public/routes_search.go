package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type searchController struct {
	service     app.SearchService
	bookService app.BookService
}

func newSearchController(service app.SearchService, bookService app.BookService) *searchController {
	return &searchController{
		service:     service,
		bookService: bookService,
	}
}

func (c *searchController) Register(r chi.Router) {
	r.Get("/search", c.search)
	r.Get("/s", c.simpleSearch)
	r.Get("/random", c.random)
}

func (c *searchController) random(w http.ResponseWriter, r *http.Request) {
	id, err := c.bookService.GetRandomBookID(r.Context())
	if err != nil {
		write500(w, r, err)
		return
	}

	if id.Valid {
		http.Redirect(w, r, fmt.Sprintf("/book/%d", id.Value), http.StatusFound)
	} else {
		http.Redirect(w, r, "/s?ol.error=no_books", http.StatusFound)
	}
}

func (c *searchController) simpleSearch(w http.ResponseWriter, r *http.Request) {
	templates.SimpleSearch().Render(r.Context(), w)
}

func (c *searchController) search(w http.ResponseWriter, r *http.Request) {
	search := getSearchRequest(r)
	query := app.BookSearchQuery{
		UserID: auth.GetNullableUserID(r.Context()),

		Query: search.Query,

		IncludeUsers:    search.IncludeUsers,
		ExcludeUsers:    search.ExcludeUsers,
		IncludeTags:     search.IncludeTags,
		ExcludeTags:     search.ExcludeTags,
		Words:           search.Words,
		Chapters:        search.Chapters,
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

	if r.URL.Query().Has("__fragment") {
		writeTemplate(w, r.Context(), templates.SearchResultFragment(result, explainedQuery))
	} else {
		writeTemplate(w, r.Context(), templates.SearchPage(result, explainedQuery))
	}
}
