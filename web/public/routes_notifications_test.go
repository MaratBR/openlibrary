package public

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	user      uuid.UUID
	after     int64
	cancel    context.CancelFunc
	through   time.Time
	statusErr error
	markErr   error
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

func (f *fixtureNotificationReader) Status(_ context.Context, user uuid.UUID) (notifications.Status, error) {
	f.user = user
	return notifications.Status{UnreadCount: 3, ObservedAt: time.Unix(1700000000, 0).UTC()}, f.statusErr
}
func (f *fixtureNotificationReader) MarkRead(_ context.Context, user uuid.UUID, through time.Time) (notifications.Status, error) {
	f.user, f.through = user, through
	return notifications.Status{LastReadAt: &through, ObservedAt: time.Unix(1700000000, 0).UTC()}, f.markErr
}
func (f *fixtureNotificationReader) Send(context.Context, string, string, string, string, ...notifications.Options) (int64, error) {
	return 0, nil
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
				} else if path == "/_api/notifications" && tc.authenticated {
					require.Equal(t, 200, w.Code)
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

func TestNotificationStatusAndReadUseAuthenticatedRecipient(t *testing.T) {
	reader := &fixtureNotificationReader{}
	router := notificationTestRouter(app.RoleUser, reader)
	for _, path := range []string{"/_api/notifications/status", "/_api/notifications/read"} {
		method, body := "GET", ""
		if strings.HasSuffix(path, "/read") {
			method, body = "POST", `{"readThrough":"2023-11-14T22:13:20Z"}`
		}
		for _, authenticated := range []bool{false, true} {
			r := httptest.NewRequest(method, path, strings.NewReader(body))
			if authenticated {
				r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			if !authenticated {
				require.Equal(t, 401, w.Code)
				continue
			}
			require.Equal(t, 200, w.Code)
			require.Equal(t, uuid.Nil(), reader.user)
			require.Equal(t, "no-store", w.Header().Get("Cache-Control"))
			if method == "POST" {
				require.Equal(t, time.Unix(1700000000, 0).UTC(), reader.through)
				require.Contains(t, w.Body.String(), `"unreadCount":0`)
			} else {
				require.Contains(t, w.Body.String(), `"unreadCount":3`)
			}
		}
	}
}

func TestNotificationReadRejectsInvalidBody(t *testing.T) {
	router := notificationTestRouter(app.RoleModerator, &fixtureNotificationReader{})
	for _, body := range []string{"", "{}", `{"readThrough":"bad"}`, `{"readThrough":null}`,
		`{"readThrough":"2023-11-14T22:13:20Z","userId":"other"}`,
		`{"readThrough":"2023-11-14T22:13:20Z"} {}`, strings.Repeat("x", 1025)} {
		r := httptest.NewRequest("POST", "/_api/notifications/read", strings.NewReader(body))
		r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		require.Equal(t, 400, w.Code, body)
	}
}

func TestNotificationStatusAndReadStorageErrors(t *testing.T) {
	for _, tc := range []struct {
		method, path, body string
		reader             *fixtureNotificationReader
		status             int
	}{
		{"GET", "/_api/notifications/status", "", &fixtureNotificationReader{statusErr: errors.New("db unavailable")}, 500},
		{"POST", "/_api/notifications/read", `{"readThrough":"2023-11-14T22:13:20Z"}`, &fixtureNotificationReader{markErr: errors.New("db unavailable")}, 500},
		{"POST", "/_api/notifications/read", `{"readThrough":"2023-11-14T22:13:20Z"}`, &fixtureNotificationReader{markErr: notifications.ErrInvalidReadTime}, 400},
	} {
		r := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
		w := httptest.NewRecorder()
		notificationTestRouter(app.RoleAdmin, tc.reader).ServeHTTP(w, r)
		require.Equal(t, tc.status, w.Code)
	}
}
