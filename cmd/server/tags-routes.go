package main

import (
	"net/http"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
)

type tagsController struct {
	tagsService *app.TagsService
}

func newTagsController(service *app.TagsService) *tagsController {
	return &tagsController{tagsService: service}
}

type searchTagsResponse struct {
	Tags  []app.DefinedTagDto `json:"tags"`
	Query string              `json:"query"`
}

func (t *tagsController) Search(w http.ResponseWriter, r *http.Request) {
	q := strings.Trim(r.URL.Query().Get("q"), " \n\t")
	if len(q) > 50 {
		q = q[:50]
	}

	tags, err := t.tagsService.SearchTags(r.Context(), q)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeJSON(w, searchTagsResponse{
		Tags:  tags,
		Query: q,
	})
}
