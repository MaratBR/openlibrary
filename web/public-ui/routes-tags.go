package publicui

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/go-chi/chi/v5"
)

type tagsController struct {
	service app.TagsService
}

func newTagsController(service app.TagsService) *tagsController {
	return &tagsController{service: service}
}

func (t *tagsController) TagPage(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "tagID")
	if tagID == "" {
		http.Redirect(w, r, "/tags", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/search?it="+tagID, http.StatusFound)
}
