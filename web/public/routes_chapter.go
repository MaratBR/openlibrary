package public

import (
	"log/slog"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type chaptersController struct {
	service            app.BookService
	readingListService app.ReadingListService
	commentsService    app.CommentsService
	viewsService       analytics.ViewsService
}

func newChaptersController(service app.BookService, readingListService app.ReadingListService, viewsService analytics.ViewsService, commentsService app.CommentsService) *chaptersController {
	return &chaptersController{service: service, readingListService: readingListService, viewsService: viewsService, commentsService: commentsService}
}

func (c *chaptersController) Register(r chi.Router) {
	r.Get("/book/{bookID}/chapters/{chapterID}", c.chapter)
	r.Get("/book/{bookID}/chapters/{chapterID}/comments", c.chapterComments)
}

func (c *chaptersController) chapter(w http.ResponseWriter, r *http.Request) {
	rl := r.URL.Query().Get("rl")

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

		status, err := c.readingListService.GetStatus(r.Context(), session.UserID, bookID)
		if err == nil && status.Valid && status.Value.ChapterID.Valid {
			statusChapterOrder := status.Value.ChapterOrder
			chapterOrder := result.Chapter.Order
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

		// if rl is 1 then update status
		if rl == "1" {
			readingListStatus, err := c.readingListService.GetStatus(r.Context(), session.UserID, result.Chapter.BookID)
			if err != nil {
				slog.Warn("readingListService.GetStatus error", "err", err)
			} else {
				if !readingListStatus.Valid || readingListStatus.Value.Status == app.ReadingListStatusWantToRead || readingListStatus.Value.Status == app.ReadingListStatusReading {
					err = c.readingListService.MarkAsReadingWithChapterID(
						r.Context(),
						session.UserID,
						book.ID,
						chapterID,
					)
					if err != nil {
						slog.Warn("readingListService.MarkAsReadingWithChapterID error", "err", err)
					}
				}
			}
		}
	}

	templates.Chapter(result.Chapter, book, options).Render(r.Context(), w)
}

func (c *chaptersController) chapterComments(w http.ResponseWriter, r *http.Request) {
	chapterID, err := olhttp.URLParamInt64(r, "chapterID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	cursor, _ := olhttp.URLQueryParamInt64(r, "cursor")

	result, err := c.commentsService.GetList(r.Context(), app.GetCommentsQuery{
		ChapterID:   chapterID,
		ActorUserID: auth.GetNullableUserID(r.Context()),
		Limit:       30,
		Cursor:      uint32(cursor),
	})

	olhttp.WriteTemplate(w, r.Context(), templates.ChapterComments(result.Comments))
}
