package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/web/olresponse"
)

type apiTagsController struct {
	service app.TagsService
}

func newAPITagsController(service app.TagsService) *apiTagsController {
	return &apiTagsController{service: service}
}

func (t *apiTagsController) Tags(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tags, err := t.service.SearchTags(r.Context(), query)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olresponse.NewAPIResponse(tags).Write(w)
}
