package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type bookManagerController struct {
	service           app.BookManagerService
	collectionService app.CollectionService
}

func newBookManagerController(service app.BookManagerService, collectionService app.CollectionService) *bookManagerController {
	return &bookManagerController{service: service, collectionService: collectionService}
}

func (c *bookManagerController) Register(r chi.Router) {
	r.Route("/books-manager", func(r chi.Router) {
		r.Use(redirectToLoginOnUnauthorized)

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

	session := auth.RequireSession(r.Context())

	if tab == "" {
		templates.BookManager().Render(r.Context(), w)
		return
	} else if tab == "books" {
		page, _ := olhttp.URLQueryParamInt64(r, "p")
		if page < 1 {
			page = 1
		} else if page > 10000 {
			page = 10000
		}
		books, err := c.service.GetUserBooks(r.Context(), app.GetUserBooksQuery{
			UserID:   session.UserID,
			PageSize: 20,
			Page:     uint32(page),
		})
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}

		templates.BookManagerBooks(books).Render(r.Context(), w)
		return
	} else if tab == "collections" {
		page := olhttp.GetPage(r.URL.Query(), "page")
		pageSize := olhttp.GetPageSize(r.URL.Query(), "pageSize", 1, 100, 15)
		result, err := c.collectionService.GetUserCollections(r.Context(), app.GetUserCollectionsQuery{
			UserID:   session.UserID,
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}

		templates.BookManagerCollections(result.Collections).Render(r.Context(), w)
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
	Rating string `in:"form=ageRating"`
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
		if err == store.ErrNoRows {
			http.Redirect(w, r, "/books-manager?tab=books", http.StatusFound)
		}

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

	session := auth.RequireSession(r.Context())

	var draftID int64

	{
		draftIDNullable, err := c.service.GetLatestDraft(r.Context(), app.GetLatestDraftQuery{ChapterID: chapterID, UserID: session.UserID})
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}
		if draftIDNullable.Valid {
			draftID = draftIDNullable.Value
		}
	}

	if draftID == 0 {
		newDraftID, err := c.service.CreateDraft(r.Context(), app.CreateDraftCommand{
			ChapterID: chapterID,
			UserID:    session.UserID,
		})

		if err != nil {
			writeApplicationError(w, r, err)
			return
		}

		draftID = newDraftID
	}

	c.sendChapterEditorPage(bookID, chapterID, draftID, w, r)
}

func (c *bookManagerController) sendChapterEditorPage(bookID, chapterID, draftID int64, w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	// chapter, err := c.service.GetChapter(r.Context(), app.ManagerGetChapterQuery{
	// 	UserID:    session.UserID,
	// 	BookID:    bookID,
	// 	ChapterID: chapterID,
	// })

	// if err != nil {
	// 	writeApplicationError(w, r, err)
	// 	return
	// }

	draft, err := c.service.GetDraft(r.Context(), app.GetDraftQuery{
		DraftID:   draftID,
		ChapterID: chapterID,
		BookID:    bookID,
		UserID:    session.UserID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.ChapterEditor(bookID, draft).Render(r.Context(), w)
}
