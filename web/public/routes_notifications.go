package public

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app/notifications"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/go-chi/chi/v5"
)

type notificationReader interface {
	History(context.Context, uuid.UUID, int64) ([]notifications.Event, error)
}

type notificationsController struct {
	reader notificationReader
}

func newNotificationsController(reader *notifications.StreamReader) *notificationsController {
	return &notificationsController{reader}
}

func (c *notificationsController) Register(r chi.Router) {
	r.Get("/_api/notifications", c.history)
}

func (c *notificationsController) history(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if !authorizeEvents(w, r) {
		return
	}
	before := int64(0)
	var err error
	if value := r.URL.Query().Get("before"); value != "" {
		before, err = strconv.ParseInt(value, 10, 64)
		if before < 0 {
			err = fmt.Errorf("negative cursor")
		}
	}
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	events, err := c.reader.History(r.Context(), auth.RequireUser(r.Context()).ID, before)
	if err != nil {
		apiWrite500(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
