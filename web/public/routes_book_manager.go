package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	i18nProvider "github.com/MaratBR/openlibrary/internal/i18n-provider"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type bookManagerController struct {
	service app.BookManagerService
}

func newBookManagerController(service app.BookManagerService) *bookManagerController {
	return &bookManagerController{service: service}
}

func (c *bookManagerController) Register(r chi.Router) {
	r.Route("/books-manager", func(r chi.Router) {
		r.Get("/", c.index)
		r.Get("/new", c.newBook)
		r.With(httpin.NewInput(&createBookRequest{})).Post("/new", c.createBook)
		r.Get("/book/{bookID}", c.book)
		r.With(httpin.NewInput(&updateBookRequest{})).Post("/book/{bookID}", c.updateBook)

	})
}

func (c *bookManagerController) index(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	tab := query.Get("tab")

	session, ok := auth.GetSession(r.Context())

	if !ok {
		redirectToLogin(w, r)
		return
	}

	if tab == "" {
		templates.BookManager().Render(r.Context(), w)
		return
	} else if tab == "books" {
		books, err := c.service.GetUserBooks(r.Context(), app.GetUserBooksQuery{
			UserID: session.UserID,
			Limit:  50,
			Offset: 0,
		})
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}

		templates.BookManagerBooks(books).Render(r.Context(), w)
		return
	} else if tab == "collections" {
		templates.BookManagerCollections().Render(r.Context(), w)
		return
	} else {
		http.Redirect(w, r, "/books-manager", http.StatusFound)
	}
}

func (c *bookManagerController) newBook(w http.ResponseWriter, r *http.Request) {
	templates.BookManagerNewBook().Render(r.Context(), w)
}

type createBookRequest struct {
	Name   string `in:"form=name"`
	Tags   string `in:"form=tags"`
	Rating string `in:"form=rating"`
}

func (c *bookManagerController) createBook(w http.ResponseWriter, r *http.Request) {
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

type updateBookRequest struct {
	Name    string `in:"form=name"`
	Summary string `in:"form=summary"`
	Tags    string `in:"form=tags"`
	Rating  string `in:"form=rating"`
}

func (c *bookManagerController) updateBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	input := r.Context().Value(httpin.Input).(*updateBookRequest)

	rating := app.AsRating(input.Rating)
	tags := olhttp.ParseInt64Array(input.Tags)
	name := input.Name
	summary := input.Summary

	err = c.service.UpdateBook(r.Context(), app.UpdateBookCommand{
		BookID:    bookID,
		Tags:      tags,
		Name:      name,
		Summary:   summary,
		AgeRating: rating,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	l := i18nProvider.GetLocalizer(r.Context())

	flash.Add(r, flash.Text(l.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "bookManager.edit.editedSuccessfully",
		TemplateData: map[string]string{
			"Name": name,
		},
	})))

	c.book(w, r)
}

func (c *bookManagerController) book(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	c.sendBookEditorPage(bookID, w, r)
}

func (c *bookManagerController) sendBookEditorPage(bookID int64, w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	query := app.ManagerGetBookQuery{
		BookID:      bookID,
		ActorUserID: session.UserID,
	}

	book, err := c.service.GetBook(r.Context(), query)

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.BookManagerBook(&book.Book).Render(r.Context(), w)
}
