package public

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/events"
	"github.com/MaratBR/openlibrary/internal/app/notifications"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type fixtureNotificationReader struct {
	user   uuid.UUID
	after  int64
	cancel context.CancelFunc
}

func (f *fixtureNotificationReader) Read(_ context.Context, user uuid.UUID, after int64) ([]notifications.Event, error) {
	f.user, f.after = user, after
	if f.cancel != nil {
		f.cancel()
	}
	return []notifications.Event{{ID: 42, Title: "Hello", Content: "<script>plain text</script>"}}, nil
}
func (f *fixtureNotificationReader) History(ctx context.Context, user uuid.UUID, before int64) ([]notifications.Event, error) {
	return f.Read(ctx, user, before)
}

func notificationTestRouter(role app.UserRole, reader *fixtureNotificationReader) http.Handler {
	log := zap.NewNop().Sugar()
	users := overridePromptUsers{role: role}
	sessions := overridePromptSessions{}
	router := chi.NewRouter()
	router.Use(olhttp.ReqCtxMiddleware, session.Middleware(overridePromptStore{}, log), auth.NewAuthorizationMiddleware(sessions, users, auth.MiddlewareOptions{}, log))
	(&notificationsController{reader: reader}).Register(router)
	(&eventsController{hub: events.NewHub(), sources: []events.Source{reader}, sessions: sessions, users: users}).Register(router)
	return router
}
func TestNotificationAccessAndCapability(t *testing.T) {
	for _, tc := range []struct {
		role                   app.UserRole
		authenticated, allowed bool
		status                 int
	}{
		{app.RoleUser, false, false, 401}, {app.RoleUser, true, false, 403}, {app.RoleModerator, true, true, 200}, {app.RoleAdmin, true, true, 200}, {app.RoleSystem, true, true, 200},
	} {
		t.Run(string(tc.role)+http.StatusText(tc.status), func(t *testing.T) {
			router := notificationTestRouter(tc.role, &fixtureNotificationReader{})
			paths := []string{"/_api/sse/capabilities", "/_api/notifications"}
			if !tc.allowed {
				paths = append(paths, "/_api/sse/stream")
			}
			for _, path := range paths {
				r := httptest.NewRequest("GET", path, nil)
				if tc.authenticated {
					r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
				}
				w := httptest.NewRecorder()
				router.ServeHTTP(w, r)
				if path == "/_api/sse/capabilities" {
					require.Equal(t, 200, w.Code)
					if tc.allowed {
						require.JSONEq(t, `{"sseAllowed":true,"userId":"00000000-0000-0000-0000-000000000000"}`, w.Body.String())
					} else {
						require.JSONEq(t, `{"sseAllowed":false}`, w.Body.String())
					}
				} else {
					require.Equal(t, tc.status, w.Code)
				}
			}
		})
	}
}
func TestNotificationStreamCursorAndDisconnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reader := &fixtureNotificationReader{cancel: cancel}
	router := notificationTestRouter(app.RoleModerator, reader)
	r := httptest.NewRequest("GET", "/_api/sse/stream?after=1", nil).WithContext(ctx)
	r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
	r.Header.Set("Last-Event-ID", encodeEventCursor(map[string]int64{"notifications": 41}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
	require.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
	require.EqualValues(t, 41, reader.after)
	require.Contains(t, w.Body.String(), "id: "+encodeEventCursor(map[string]int64{"notifications": 42})+"\nevent: user-event\ndata: ")
	require.Contains(t, w.Body.String(), "\"title\":\"Hello\"")
	require.Contains(t, w.Body.String(), ": heartbeat")
}
func TestNotificationInvalidCursors(t *testing.T) {
	router := notificationTestRouter(app.RoleAdmin, nil)
	for _, path := range []string{"/_api/sse/stream?after=-1", "/_api/sse/stream?after=bad", "/_api/notifications?before=-1", "/_api/notifications?before=9223372036854775808"} {
		r := httptest.NewRequest("GET", path, nil)
		r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		require.Equal(t, 400, w.Code)
	}
}

func (f *fixtureNotificationReader) Name() string { return "notifications" }
func (f *fixtureNotificationReader) Poll(ctx context.Context, user uuid.UUID, after int64) ([]events.StoredEvent, error) {
	records, err := f.Read(ctx, user, after)
	if err != nil {
		return nil, err
	}
	result := make([]events.StoredEvent, 0, len(records))
	for _, record := range records {
		data, _ := json.Marshal(record)
		result = append(result, events.StoredEvent{Cursor: record.ID, Event: events.Event{Type: "notification", Data: data}})
	}
	return result, nil
}
