package public

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/events"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type eventsController struct {
	hub      events.Hub
	sources  []events.Source
	sessions app.SessionService
	users    app.UserService
}

type eventsControllerParams struct {
	fx.In
	Hub      events.Hub
	Sources  []events.Source `group:"user_sse_sources"`
	Sessions app.SessionService
	Users    app.UserService
}

func newEventsController(p eventsControllerParams) *eventsController {
	return &eventsController{hub: p.Hub, sources: p.Sources, sessions: p.Sessions, users: p.Users}
}

// Resume state is opaque to clients and keeps each durable source independent.
func eventCursor(r *http.Request) (map[string]int64, error) {
	value := r.Header.Get("Last-Event-ID")
	if value == "" {
		value = r.URL.Query().Get("after")
	}
	cursors := make(map[string]int64)
	if value == "" || value == "0" {
		return cursors, nil
	}
	if len(value) > 4096 {
		return nil, fmt.Errorf("invalid cursor")
	}
	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor")
	}
	if err := json.Unmarshal(data, &cursors); err != nil || cursors == nil {
		return nil, fmt.Errorf("invalid cursor")
	}
	for _, cursor := range cursors {
		if cursor < 0 {
			return nil, fmt.Errorf("invalid cursor")
		}
	}
	return cursors, nil
}
func encodeEventCursor(cursors map[string]int64) string {
	data, _ := json.Marshal(cursors)
	return base64.RawURLEncoding.EncodeToString(data)
}

func (c *eventsController) Register(r chi.Router) {
	r.Get("/_api/sse/capabilities", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		allowed := eventsAllowed(r)
		userID := ""
		if allowed {
			userID = auth.RequireUser(r.Context()).ID.String()
		}
		json.NewEncoder(w).Encode(struct {
			SSEAllowed bool   `json:"sseAllowed"`
			UserID     string `json:"userId,omitempty"`
		}{allowed, userID})
	})
	r.Get("/_api/sse/stream", c.stream)
}

func mustEventJSON(event events.Event) []byte {
	// Both live publishing and durable adapters supply valid JSON.
	data, _ := json.Marshal(event)
	return data
}

func (c *eventsController) stream(w http.ResponseWriter, r *http.Request) {
	if !authorizeEvents(w, r) {
		return
	}
	cursors, err := eventCursor(r)
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	controller := http.NewResponseController(w)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	if err := setEventWriteDeadline(controller); err != nil {
		return
	}
	if _, err := fmt.Fprint(w, "retry: 3000\n: connected\n\n"); err != nil {
		return
	}
	if err := controller.Flush(); err != nil {
		return
	}
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	// Periodically reconnect through the normal auth middleware and bound request lifetime.
	lifetime := time.NewTimer(45 * time.Second)
	defer lifetime.Stop()
	messages, unsubscribe := c.hub.Subscribe(auth.RequireUser(r.Context()).ID)
	defer unsubscribe()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-lifetime.C:
			return
		default:
		}
		user := c.streamUser(r)
		if user == nil {
			return
		}

		hasMore := false
		if err := setEventWriteDeadline(controller); err != nil {
			return
		}
		for _, source := range c.sources {
			batch, err := source.Poll(r.Context(), user.ID, cursors[source.Name()])
			if err != nil {
				return
			}
			for _, record := range batch {
				cursors[source.Name()] = record.Cursor
				if _, err := fmt.Fprintf(w, "id: %s\nevent: user-event\ndata: %s\n\n", encodeEventCursor(cursors), mustEventJSON(record.Event)); err != nil {
					return
				}
			}
			if len(batch) == 100 {
				hasMore = true
			}
		}
		if _, err := fmt.Fprint(w, ": heartbeat\n\n"); err != nil {
			return
		}
		if err := controller.Flush(); err != nil {
			return
		}
		if hasMore {
			continue
		}
		select {
		case <-r.Context().Done():
			return
		case <-lifetime.C:
			return
		case message, ok := <-messages:
			if !ok {
				return
			}
			// Recheck authorization before emitting a live event.
			if c.streamUser(r) == nil {
				return
			}
			if err := setEventWriteDeadline(controller); err != nil {
				return
			}
			if _, err := fmt.Fprintf(w, "event: user-event\ndata: %s\n\n", mustEventJSON(message)); err != nil {
				return
			}
			if err := controller.Flush(); err != nil {
				return
			}
		case <-ticker.C:
		}
	}
}
func setEventWriteDeadline(controller *http.ResponseController) error {
	err := controller.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if errors.Is(err, http.ErrNotSupported) {
		return nil
	}
	return err
}

func (c *eventsController) streamUser(r *http.Request) *app.SelfUserDto {
	session := auth.RequireSession(r.Context())
	current, err := c.sessions.GetBySID(r.Context(), session.SID)
	if err != nil {
		return nil
	}
	user, err := c.users.GetUserSelfData(r.Context(), current.UserID)
	if err != nil || !user.Role.IsModOrHigher() {
		return nil
	}
	return user
}

func eventsAllowed(r *http.Request) bool {
	user, ok := auth.GetUser(r.Context())
	_, sessionOK := auth.GetSession(r.Context())
	return ok && sessionOK && user.Role.IsModOrHigher()
}

func authorizeEvents(w http.ResponseWriter, r *http.Request) bool {
	if _, ok := auth.GetSession(r.Context()); !ok {
		apiWriteUnauthorized(w)
		return false
	}
	if !eventsAllowed(r) {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}
