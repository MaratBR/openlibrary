package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type tagsController struct {
	service app.TagsService
}

func newTagsController(service app.TagsService) *tagsController {
	return &tagsController{service: service}
}

func (t *tagsController) Register(r chi.Router) {
	r.Get("/tag/{tagID}", t.TagPage)
}

func (t *tagsController) TagPage(w http.ResponseWriter, r *http.Request) {
	tagID, err := olhttp.URLParamInt64(r, "tagID")
	if err != nil {
		http.Redirect(w, r, "/tags", http.StatusFound)
		return
	}

	tag, err := t.service.GetTag(r.Context(), tagID)

	if err != nil {
		if apperror.IsNotFoundError(err) {
			http.Redirect(w, r, "/tags", http.StatusFound)
			return
		}

		writeApplicationError(w, r, err)
		return
	}

	olhttp.WriteTemplate(w, r.Context(), templates.TagPage(tag))
}
