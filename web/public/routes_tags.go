package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
)

type tagsController struct {
	service app.TagsService
}

func newTagsController(service app.TagsService) *tagsController {
	return &tagsController{service: service}
}

func (t *tagsController) TagPage(w http.ResponseWriter, r *http.Request) {
	tagID, err := olhttp.URLParamInt64(r, "tagID")
	if err != nil {
		http.Redirect(w, r, "/tags", http.StatusFound)
		return
	}

	tag, err := t.service.GetTag(r.Context(), tagID)

	if err != nil {
		if app.IsNotFoundError(err) {
			http.Redirect(w, r, "/tags", http.StatusFound)
			return
		}

		writeApplicationError(w, r, err)
		return
	}

	writeTemplate(w, r.Context(), templates.TagPage(tag))
}
