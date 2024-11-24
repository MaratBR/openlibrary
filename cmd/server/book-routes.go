package server

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/store"
)

type bookController struct {
	bookService *app.BookService
}

func newBookController(bookService *app.BookService) *bookController {
	return &bookController{
		bookService: bookService,
	}
}

type myBooksResponse struct {
	Books []app.ManagerAuthorBookDto `json:"books"`
}

func (c *bookController) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "id")
	if err != nil {
		writeRequestError(err, w)
		return
	}
	if bookID == 0 {
		writeUnprocessableEntity(w, "book id must be provided")
		return
	}

	nullableUserID := getNullableUserID(r)
	if err == nil {
	}
	book, err := c.bookService.GetBook(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: nullableUserID})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, book)
}

type getChapterResponse struct {
	Chapter app.ChapterDto `json:"chapter"`
}

func (c *bookController) GetChapter(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeRequestError(err, w)
		return
	}
	chapterID, err := urlParamInt64(r, "chapterID")
	if err != nil {
		writeRequestError(err, w)
		return
	}
	chapter, err := c.bookService.GetBookChapter(r.Context(), app.GetBookChapterQuery{
		BookID:    bookID,
		ChapterID: chapterID,
	})
	if err != nil {
		if err == store.ErrNoRows {
			write404(w, "chapter not found")
			return
		}
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, getChapterResponse{
		Chapter: chapter.Chapter.Chapter,
	})
}
