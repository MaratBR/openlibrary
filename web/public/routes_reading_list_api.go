package public

import (
	"context"
	"errors"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/olresponse"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type apiReadingListController struct {
	service app.ReadingListService
}

func newAPIReadingListController(service app.ReadingListService) *apiReadingListController {
	return &apiReadingListController{service: service}
}

func (c *apiReadingListController) Register(r chi.Router) {
	r.Post("/reading-list/status", c.UpdateStatus)
	r.Post("/reading-list/chapter", c.UpdateCurrentChapter)
}

func (c *apiReadingListController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	bookID, err := olhttp.URLQueryParamInt64(r, "bookId")
	if err != nil {
		apiWriteUnprocessableEntity(w, err)
		return
	}
	if bookID == 0 {
		apiWriteUnprocessableEntity(w, errors.New("bookId cannot be 0"))
		return
	}

	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		apiWriteUnprocessableEntity(w, errors.New("status query parameter is missing"))
		return
	}

	var f func(ctx context.Context, userID uuid.UUID, bookID int64) error

	switch statusStr {
	case string(app.ReadingListStatusDnf):
		f = c.service.MarkAsDnf
	case string(app.ReadingListStatusRead):
		f = c.service.MarkAsRead
	case string(app.ReadingListStatusWantToRead):
		f = c.service.MarksAsWantToRead
	case string(app.ReadingListStatusPaused):
		f = c.service.MarkAsPaused
	case string(app.ReadingListStatusReading):
		f = c.service.MarkAsReading
	default:
		apiWriteUnprocessableEntity(w, errors.New("invalid value for status: "+statusStr))
		return
	}

	err = f(r.Context(), session.UserID, bookID)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	state, err := c.service.GetStatus(r.Context(), session.UserID, bookID)
	if err != nil {
		apiWriteUnexpectedApplicationError(w, err)
		return
	}

	if !state.Valid {
		apiWrite500(w, errors.New("GetStatus returned no value (null) after updating the state of the book"))
		return
	}

	olresponse.NewAPIResponse(state.Value).Write(w)
}

func (c *apiReadingListController) UpdateCurrentChapter(w http.ResponseWriter, r *http.Request) {
	chapterID, err := olhttp.URLQueryParamInt64(r, "chapterId")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	s := auth.RequireSession(r.Context())

	err = c.service.MarkChapterRead(r.Context(), app.MarkChapterCommand{
		UserID:    s.UserID,
		ChapterID: chapterID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
	} else {
		apiWriteOK(w)
	}
}
