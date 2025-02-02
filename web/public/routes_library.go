package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type libraryController struct{}

func newLibraryController() *libraryController {
	return &libraryController{}
}

func (c *libraryController) Register(r chi.Router) {
	r.Get("/", c.index)
}

func (c *libraryController) index(w http.ResponseWriter, r *http.Request) {
	templates.Library().Render(r.Context(), w)
}
