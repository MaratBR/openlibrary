package public

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/notifications"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type notificationsController struct {
	reader notifications.Service
}

func newNotificationsController(reader notifications.Service) *notificationsController {
	return &notificationsController{reader}
}

func (c *notificationsController) Register(r chi.Router) {
	r.With(requiresAuthorizationMiddleware).Get("/notifications", c.center)
	r.Get("/_api/notifications", c.history)
	r.Get("/_api/notifications/status", c.status)
	r.Post("/_api/notifications/read", c.markRead)
}

func (c *notificationsController) history(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if !authorizeNotificationInbox(w, r) {
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

func authorizeNotificationInbox(w http.ResponseWriter, r *http.Request) bool {
	_, sessionOK := auth.GetSession(r.Context())
	_, userOK := auth.GetUser(r.Context())
	if !sessionOK || !userOK {
		apiWriteUnauthorized(w)
		return false
	}
	w.Header().Set("Cache-Control", "no-store")
	return true
}

func (c *notificationsController) status(w http.ResponseWriter, r *http.Request) {
	if !authorizeNotificationInbox(w, r) {
		return
	}
	status, err := c.reader.Status(r.Context(), auth.RequireUser(r.Context()).ID)
	if err != nil {
		apiWrite500(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (c *notificationsController) markRead(w http.ResponseWriter, r *http.Request) {
	if !authorizeNotificationInbox(w, r) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var input struct {
		ReadThrough time.Time `json:"readThrough"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		apiWriteBadRequest(w, errors.New("unexpected JSON content"))
		return
	}
	if input.ReadThrough.IsZero() {
		apiWriteBadRequest(w, notifications.ErrInvalidReadTime)
		return
	}
	status, err := c.reader.MarkRead(r.Context(), auth.RequireUser(r.Context()).ID, input.ReadThrough)
	if err != nil {
		if errors.Is(err, notifications.ErrInvalidReadTime) {
			apiWriteBadRequest(w, err)
		} else {
			apiWrite500(w, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (c *notificationsController) center(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	templates.NotificationsCenter().Render(r.Context(), w)
}
