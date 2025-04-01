package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
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
		r.Get("/book/{bookID}", c.book)
		r.Get("/book/{bookID}/chapter/{chapterID}", c.chapter)

		r.With(httpin.NewInput(&createBookRequest{})).Post("/new", c.createBook)

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

func (c *bookManagerController) chapter(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	chapterID, err := olhttp.URLParamInt64(r, "chapterID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	c.sendChapterEditorPage(bookID, chapterID, w, r)
}

func (c *bookManagerController) sendChapterEditorPage(bookID, chapterID int64, w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	chapter, err := c.service.GetChapter(r.Context(), app.ManagerGetChapterQuery{
		UserID:    session.UserID,
		BookID:    bookID,
		ChapterID: chapterID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.ChapterEditor(chapter.Chapter).Render(r.Context(), w)
}
