package admin

import (
	"errors"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/notifications"
	"go.uber.org/zap"
)

type debugActionDescriptor struct {
	Func func() string
}

type debugController struct {
	fullReindexService app.BookReindexService
	notifications      *notifications.StreamReader
	actions            map[string]debugActionDescriptor
}

func newDebugController(fullReindexService app.BookReindexService, notifications *notifications.StreamReader) *debugController {
	c := &debugController{
		fullReindexService: fullReindexService,
		notifications:      notifications,
		actions: map[string]debugActionDescriptor{
			"books:elastic:reindex": {
				Func: func() string {
					err := fullReindexService.ScheduleReindexAll()
					if err == nil {
						return ""
					}

					return err.Error()
				},
			},
		},
	}
	return c
}

func (c *debugController) sendNotification(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if err := r.ParseForm(); err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	count, err := c.notifications.Send(r.Context(), r.FormValue("userId"), r.FormValue("role"), r.FormValue("title"), r.FormValue("content"))
	if err != nil {
		if errors.Is(err, notifications.ErrInvalidSend) {
			adminJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			zap.S().Errorw("failed to persist custom notifications", "error", err)
			adminAPIError(w, http.StatusInternalServerError)
		}
		return
	}
	adminJSON(w, http.StatusOK, struct {
		Count int64 `json:"count"`
	}{count})
}
