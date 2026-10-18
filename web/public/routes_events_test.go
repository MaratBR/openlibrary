package public

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/events"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type liveEventWriter struct {
	*httptest.ResponseRecorder
	hub     events.Hub
	cancel  context.CancelFunc
	flushes int
}

func (w *liveEventWriter) Flush() {
	w.ResponseRecorder.Flush()
	w.flushes++
	if w.flushes == 2 {
		_ = w.hub.Send(uuid.UUID{}, "book.progress", map[string]int{"progress": 50})
	}
	if w.flushes == 3 {
		w.cancel()
	}
}
func TestSSETransportsArbitraryUserEventsWithoutNotificationSource(t *testing.T) {
	log := zap.NewNop().Sugar()
	users := overridePromptUsers{role: app.RoleModerator}
	sessions := overridePromptSessions{}
	hub := events.NewHub()
	router := chi.NewRouter()
	router.Use(olhttp.ReqCtxMiddleware, session.Middleware(overridePromptStore{}, log), auth.NewAuthorizationMiddleware(sessions, users, auth.MiddlewareOptions{}, log))
	(&eventsController{hub: hub, sessions: sessions, users: users}).Register(router)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := httptest.NewRequest("GET", "/_api/sse/stream", nil).WithContext(ctx)
	r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
	w := &liveEventWriter{ResponseRecorder: httptest.NewRecorder(), hub: hub, cancel: cancel}
	router.ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), `data: {"type":"book.progress","data":{"progress":50}}`)
}
func TestSSECursorKeepsSourceNamespacesIndependent(t *testing.T) {
	cursors := map[string]int64{"notifications": 42, "jobs": 7}
	r := httptest.NewRequest("GET", "/_api/sse/stream", nil)
	r.Header.Set("Last-Event-ID", encodeEventCursor(cursors))
	decoded, err := eventCursor(r)
	require.NoError(t, err)
	require.Equal(t, cursors, decoded)
	for _, raw := range []string{"bad", "-1", "bnVsbA", "eyJqb2JzIjotMX0"} {
		r.Header.Set("Last-Event-ID", raw)
		_, err := eventCursor(r)
		require.Error(t, err)
	}
}
