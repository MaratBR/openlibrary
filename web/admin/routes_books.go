package admin

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/go-chi/chi/v5"
)

type booksController struct {
	db store.DBTX
}

func newBooksController(db store.DBTX) *booksController {
	return &booksController{db: db}
}

func (c *booksController) Register(r chi.Router) {
	r.Get("/", c.mainPage)
}

func (c *booksController) mainPage(w http.ResponseWriter, r *http.Request) {
	olhttp.WriteTemplate(w, r.Context(), templates.BooksPage())
}
