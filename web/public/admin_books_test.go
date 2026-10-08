package public

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/web/frontend"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"uuid"
)

type overridePromptBookService struct {
	app.BookService
	err   error
	calls int
}

func (s *overridePromptBookService) GetBookDetails(context.Context, app.GetBookQuery) (app.BookDetailsDto, error) {
	s.calls++
	return app.BookDetailsDto{}, s.err
}

type overridePromptSession struct{}

func (overridePromptSession) Get(string) (string, bool)  { return "", false }
func (overridePromptSession) Put(string, string)         {}
func (overridePromptSession) Save(context.Context) error { return nil }
func (overridePromptSession) ID() string                 { return "fixture" }

type overridePromptStore struct{}

func (overridePromptStore) Get(context.Context, string) (session.Session, error) {
	return overridePromptSession{}, nil
}

type overridePromptSessions struct{ app.SessionService }

func (overridePromptSessions) GetBySID(context.Context, string) (*app.SessionInfo, error) {
	return &app.SessionInfo{}, nil
}

type overridePromptUsers struct {
	app.UserService
	role app.UserRole
}

func (s overridePromptUsers) GetUserSelfData(context.Context, uuid.UUID) (*app.SelfUserDto, error) {
	return &app.SelfUserDto{Name: "Admin", Role: s.role}, nil
}

func TestRestrictedBookOffersOverrideOnlyToAdmins(t *testing.T) {
	for _, tc := range []struct {
		name                         string
		role                         app.UserRole
		authenticated, missing, want bool
	}{
		{"admin", app.RoleAdmin, true, false, true},
		{"system", app.RoleSystem, true, false, true},
		{"moderator", app.RoleModerator, true, false, false},
		{"user", app.RoleUser, true, false, false},
		{"anonymous", app.RoleAdmin, false, false, false},
		{"missing", app.RoleAdmin, true, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := &overridePromptBookService{err: app.ErrTypeBookPrivated.New("restricted")}
			if tc.missing {
				service.err = app.ErrTypeBookNotFound.New("missing")
			}
			log := zap.NewNop().Sugar()
			provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../translations/en.toml"}}, log)
			router := chi.NewRouter()
			router.Use(olhttp.ReqCtxMiddleware, provider.Middleware, session.Middleware(overridePromptStore{}, log), auth.NewAuthorizationMiddleware(overridePromptSessions{}, overridePromptUsers{role: tc.role}, auth.MiddlewareOptions{}, log))
			frontend.AttachAssetsInliningHandler(os.DirFS("../frontend/embed-assets"), "embed-assets", router)
			(&bookController{bookService: service}).Register(router)
			out := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/book/title-123?from=search", nil)
			if tc.authenticated {
				req.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
			}
			router.ServeHTTP(out, req)
			html := out.Body.String()
			has := strings.Contains(html, "Override and view book")
			if has != tc.want {
				t.Fatalf("prompt=%v want=%v: %s", has, tc.want, html)
			}
			if tc.want && !strings.Contains(html, `href="/book/title-123?admin.override=1&amp;from=search"`) {
				t.Fatal("override link lost book path or existing query")
			}
			if service.calls != 1 {
				t.Fatalf("restricted content was fetched again before opt-in: calls=%d", service.calls)
			}
		})
	}
}
