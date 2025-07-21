package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type chaptersController struct {
	service app.BookService
}

func newChaptersController(service app.BookService) *chaptersController {
	return &chaptersController{service: service}
}

func (c *chaptersController) Register(r chi.Router) {
	r.Get("/book/{bookID}/chapters/{chapterID}", c.GetChapter)

}

func (c *chaptersController) GetChapter(w http.ResponseWriter, r *http.Request) {

	bookID, err := commonutil.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	chapterID, err := commonutil.URLParamInt64(r, "chapterID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	userID := auth.GetNullableUserID(r.Context())
	book, err := c.service.GetBook(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: userID})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	chapter, err := c.service.GetBookChapter(r.Context(), app.GetBookChapterQuery{
		BookID:    bookID,
		ChapterID: chapterID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.Chapter(chapter.Chapter, book).Render(r.Context(), w)
}
