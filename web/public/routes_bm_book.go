package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/ggicci/httpin"
)

func (c *bookManagerController) bookNew(w http.ResponseWriter, r *http.Request) {
	templates.BM_NewBook().Render(r.Context(), w)
}

type createBookRequest struct {
	Name   string `in:"form=name"`
	Tags   string `in:"form=tags"`
	Rating string `in:"form=ageRating"`
}

func (c *bookManagerController) bookCreate(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	input := r.Context().Value(httpin.Input).(*createBookRequest)

	rating := app.AsRating(input.Rating)
	tags := olhttp.ParseInt64Array(input.Tags)
	name := input.Name

	command := app.CreateBookCommand{
		Name:              name,
		Tags:              tags,
		AgeRating:         rating,
		Summary:           "",
		IsPubliclyVisible: true,
		UserID:            session.UserID,
	}

	bookID, err := c.service.CreateBook(r.Context(), command)

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/books-manager#/books/%d", bookID), http.StatusFound)
}
