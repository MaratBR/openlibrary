package server

import (
	"net/http"

	"github.com/MaratBR/openlibrary/cmd/server/templates"
	"github.com/MaratBR/openlibrary/internal/app"
)

type searchUIController struct {
	searchService app.SearchService
	tagsService   app.TagsService
}

func newSearchUIController(searchService app.SearchService, tagsService app.TagsService) *searchUIController {
	return &searchUIController{
		searchService: searchService,
		tagsService:   tagsService,
	}
}

func (c *searchUIController) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	templates.SearchPage().Render(r.Context(), w)
}
