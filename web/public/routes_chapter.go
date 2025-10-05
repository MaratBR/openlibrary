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
	service            app.BookService
	readingListService app.ReadingListService
}

func newChaptersController(service app.BookService, readingListService app.ReadingListService) *chaptersController {
	return &chaptersController{service: service, readingListService: readingListService}
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
	book, err := c.service.GetBookDetails(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: userID})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	result, err := c.service.GetBookChapter(r.Context(), app.GetBookChapterQuery{
		BookID:    bookID,
		ChapterID: chapterID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	var options templates.ChapterProgressTrackerOptions

	session, ok := auth.GetSession(r.Context())
	if ok {
		options.Enable = true

		if session == nil {
			panic("wtf")
		}

		status, err := c.readingListService.GetStatus(r.Context(), session.UserID, bookID)
		if err == nil && status.Valid && status.Value.ChapterID.Valid {
			statusChapterOrder := status.Value.ChapterOrder
			chapterOrder := result.ChapterWithDetails.Chapter.Order
			if statusChapterOrder == chapterOrder {
				// if it's same chapter - no need to do anything, disable chapter auto-marking
				options.Enable = false
			} else if chapterOrder < statusChapterOrder {
				// we backtracked
				options.JumpedBackward = true
			} else if chapterOrder > statusChapterOrder+1 {
				// we jumped forward 1 or more over
				options.JumpedForward = true
			}
		}
	}

	templates.Chapter(result.ChapterWithDetails, book, options).Render(r.Context(), w)
}
