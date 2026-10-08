package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/go-chi/chi/v5"
)

type apiControllerTags struct {
	service app.TagsService
}

func newAPITagsController(service app.TagsService) *apiControllerTags {
	return &apiControllerTags{service: service}
}

func (t *apiControllerTags) Register(r chi.Router) {
	r.Get("/tags", t.Tags)
	r.Get("/books-manager/tags", t.Tags)
}

func (t *apiControllerTags) Tags(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tags, err := t.service.SearchTags(r.Context(), query)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	visible := make([]app.DefinedTagDto, 0, len(tags))
	for _, tag := range tags {
		if !app.IsTagCensored(r.Context(), tag.ID) {
			visible = append(visible, tag)
		}
	}
	olhttp.NewAPIResponse(visible).Write(w)
}
