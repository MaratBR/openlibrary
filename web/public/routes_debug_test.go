package public

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MaratBR/openlibrary/web/webinfra"
	"github.com/go-chi/chi/v5"
)

func TestLaunchTimeEndpoint(t *testing.T) {
	started := time.Date(2026, 10, 8, 21, 0, 0, 123456789, time.FixedZone("local", 7*60*60))
	handler := chi.NewRouter()
	newDebugController(&webinfra.ServerRuntime{StartedAt: started}).Register(handler)
	handler.Get("/search", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	for range 2 {
		out := httptest.NewRecorder()
		handler.ServeHTTP(out, httptest.NewRequest("GET", "/debug/launch-time", nil))
		if out.Code != http.StatusOK || out.Body.String() != "2026-10-08T14:00:00.123456789Z" {
			t.Fatalf("unexpected launch time: %d %s", out.Code, out.Body.String())
		}
		if out.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("launch time must not be cached")
		}
	}
	out := httptest.NewRecorder()
	handler.ServeHTTP(out, httptest.NewRequest("HEAD", "/debug/launch-time", nil))
	if out.Code != 200 || out.Body.Len() != 0 {
		t.Fatal("HEAD returned a body")
	}
	out = httptest.NewRecorder()
	handler.ServeHTTP(out, httptest.NewRequest("POST", "/debug/launch-time", nil))
	if out.Code != 405 || !strings.Contains(strings.Join(out.Header().Values("Allow"), ","), "GET") || !strings.Contains(strings.Join(out.Header().Values("Allow"), ","), "HEAD") {
		t.Fatal("invalid method accepted")
	}
	out = httptest.NewRecorder()
	handler.ServeHTTP(out, httptest.NewRequest("GET", "/search", nil))
	if out.Code != 204 {
		t.Fatal("other routes were intercepted")
	}
}
