package publicui

import (
	"context"
	"errors"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/gofrs/uuid"
)

type apiReadingListController struct {
	service app.ReadingListService
}

func newAPIReadingListController(service app.ReadingListService) *apiReadingListController {
	return &apiReadingListController{service: service}
}

func (c *apiReadingListController) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	bookID, err := olhttp.URLQueryParamInt64(r, "bookId")
	if err != nil {
		apiWriteUnprocessableEntity(w, "failed to parse bookId: "+err.Error())
		return
	}
	if bookID == 0 {
		apiWriteUnprocessableEntity(w, "bookId cannot be 0")
		return
	}

	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		apiWriteUnprocessableEntity(w, "status query parameter is missing")
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
		apiWriteUnprocessableEntity(w, "invalid value for status: "+statusStr)
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

	apiWriteJSON(w, &state.Value)
}
