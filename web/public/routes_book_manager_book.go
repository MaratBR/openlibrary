package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/ggicci/httpin"
)

func (c *bookManagerController) bookNew(w http.ResponseWriter, r *http.Request) {
	templates.BookManagerNewBook().Render(r.Context(), w)
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

	http.Redirect(w, r, fmt.Sprintf("/books-manager/book/%d", bookID), http.StatusFound)
}

func (c *bookManagerController) book(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	c.bookEditorPage(bookID, w, r)
}

func (c *bookManagerController) bookEditorPage(bookID int64, w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	query := app.ManagerGetBookQuery{
		BookID:      bookID,
		ActorUserID: session.UserID,
	}

	book, err := c.service.GetBook(r.Context(), query)

	if err != nil {
		if err == store.ErrNoRows {
			http.Redirect(w, r, "/books-manager/books", http.StatusFound)
		}

		writeApplicationError(w, r, err)
		return
	}

	templates.BookManagerBook(&book.Book).Render(r.Context(), w)
}

func (c *bookManagerController) bookUpdateGeneralInformation(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	rating := app.AsRating(r.FormValue("rating"))
	tags := olhttp.GetInt64Array(r.Form, "tags")
	name := r.FormValue("name")
	summary := r.FormValue("summary")
	isPubliclyVisible := olhttp.GetBool(r.Form, "isPubliclyVisible").Or(false)

	book, err := c.service.GetBook(r.Context(), app.ManagerGetBookQuery{
		ActorUserID: session.UserID,
		BookID:      bookID,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	err = c.service.UpdateBook(r.Context(), app.UpdateBookCommand{
		BookID:            bookID,
		Name:              name,
		AgeRating:         rating,
		Tags:              tags,
		Summary:           summary,
		IsPubliclyVisible: isPubliclyVisible,

		UserID: session.UserID,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	l := i18n.GetLocalizer(r.Context())
	flash.Add(r, flash.Text(l.TData("bookManager.edit.editedSuccessfully", map[string]string{
		"Name": book.Book.Name,
	})))

	http.Redirect(w, r, fmt.Sprintf("/books-manager/book/%d", bookID), http.StatusFound)
}
