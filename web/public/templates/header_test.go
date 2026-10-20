package templates

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

func TestHeaderActionsPrerenderedForHydration(t *testing.T) {
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../../translations/en.toml"}}, zap.NewNop().Sugar())
	router := chi.NewRouter()
	router.Use(provider.Middleware, session.Middleware(bookFixtureStore{}, zap.NewNop().Sugar()), auth.NewAuthorizationMiddleware(bookFixtureSessions{}, bookFixtureUsers{}, auth.MiddlewareOptions{}, zap.NewNop().Sugar()))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := navUser().Render(r.Context(), w); err != nil {
			t.Fatal(err)
		}
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
	out := httptest.NewRecorder()
	router.ServeHTTP(out, request)
	html := out.Body.String()
	for _, expected := range []string{`data-react-hydrate="true"`, `class="flex items-center gap-2"`, `HeaderPopover-trigger UserMenu-trigger`, `HeaderPopover-trigger NotificationMenu-trigger`, `href="/notifications"`} {
		if !strings.Contains(html, expected) {
			t.Errorf("missing initial header markup: %s", expected)
		}
	}
	for _, unexpected := range []string{`<noscript>`, `x-show="_loading"`} {
		if strings.Contains(html, unexpected) {
			t.Errorf("header actions must remain visible before hydration: %s", unexpected)
		}
	}
	if path := os.Getenv("HEADER_FIXTURE_PATH"); path != "" {
		if err := os.WriteFile(path, []byte(html), 0600); err != nil {
			t.Fatal(err)
		}
	}
}
