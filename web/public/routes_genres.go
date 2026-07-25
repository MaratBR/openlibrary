package public

import (
	"net/http"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

const genresPageSize = 24

type genresController struct {
	service app.TagsService
}

func newGenresController(service app.TagsService) *genresController {
	return &genresController{service: service}
}

func (c *genresController) Register(r chi.Router) {
	r.Get("/browse/genres", c.index)
}

func (c *genresController) index(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	result, err := c.service.List(r.Context(), app.ListTagsQuery{
		SearchQuery:    query,
		Page:           olhttp.GetPage(r.URL.Query(), "p"),
		PageSize:       genresPageSize,
		OnlyParentTags: true,
		Category:       app.Value(app.TagsCategoryGenre),
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	olhttp.WriteTemplate(w, r.Context(), templates.GenresPage(result, query))
}
