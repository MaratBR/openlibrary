package main

import (
	"context"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/gofrs/uuid"
)

type readingListController struct {
	service app.ReadingListService
}

func newReadingListController(service app.ReadingListService) *readingListController {
	return &readingListController{service: service}
}

func (c *readingListController) GetStatus(w http.ResponseWriter, r *http.Request) {
	session := requireSession(r)

	bookID, err := urlQueryParamInt64(r, "bookId")
	if err != nil {
		writeUnprocessableEntity(w, "failed to parse bookId: "+err.Error())
		return
	}
	if bookID == 0 {
		writeUnprocessableEntity(w, "bookId cannot be 0")
		return
	}

	state, err := c.service.GetStatus(r.Context(), session.UserID, bookID)
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, state)
}

func (c *readingListController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	session := requireSession(r)

	bookID, err := urlQueryParamInt64(r, "bookId")
	if err != nil {
		writeUnprocessableEntity(w, "failed to parse bookId: "+err.Error())
		return
	}
	if bookID == 0 {
		writeUnprocessableEntity(w, "bookId cannot be 0")
		return
	}

	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		writeUnprocessableEntity(w, "status query parameter is missing")
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
		writeUnprocessableEntity(w, "invalid value for status: "+statusStr)
		return
	}

	err = f(r.Context(), session.UserID, bookID)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	state, err := c.service.GetStatus(r.Context(), session.UserID, bookID)
	if err != nil {
		writeUnexpectedApplicationError(w, err)
		return
	}

	if !state.Valid {
		write500(w, "GetStatus returned no value (null) after updating the state of the book")
		return
	}

	writeJSON(w, &state.Value)
}

func (c *readingListController) Delete(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlQueryParamInt64(r, "bookId")
	if err != nil {
		writeUnprocessableEntity(w, "failed to parse bookId: "+err.Error())
		return
	}
	if bookID == 0 {
		writeUnprocessableEntity(w, "bookId cannot be 0")
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented yet"))
}

func (c *readingListController) StartReading(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlQueryParamInt64(r, "bookId")
	if err != nil {
		writeUnprocessableEntity(w, "failed to parse bookId: "+err.Error())
		return
	}
	if bookID == 0 {
		writeUnprocessableEntity(w, "bookId cannot be 0")
		return
	}

	chapterID, err := urlQueryParamInt64(r, "chapterId")
	if err != nil {
		writeUnprocessableEntity(w, "failed to parse chapterId: "+err.Error())
		return
	}
	if chapterID == 0 {
		writeUnprocessableEntity(w, "chapterId cannot be 0")
		return
	}

	session := requireSession(r)
	err = c.service.MarkAsReadingWithChapterID(r.Context(), session.UserID, bookID, chapterID)
	if err != nil {
		writeApplicationError(w, err)
	}

	state, err := c.service.GetStatus(r.Context(), session.UserID, bookID)
	if err != nil {
		writeUnexpectedApplicationError(w, err)
		return
	}

	if !state.Valid {
		write500(w, "GetStatus returned no value (null) after updating the state of the book")
		return
	}

	writeJSON(w, &state.Value)
}
