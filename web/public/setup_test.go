package public

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/siteconfig"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/frontend"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type setupPageService struct {
	state              app.GetSetupDataResult
	checkErr, setupErr error
	calls              int
	command            app.SetupCommand
}

func (s *setupPageService) IsSetup() (bool, error) { return s.state.IsSetup, s.checkErr }
func (s *setupPageService) GetSetupState() (app.GetSetupDataResult, error) {
	return s.state, s.checkErr
}
func (s *setupPageService) Setup(_ context.Context, c app.SetupCommand) error {
	s.calls++
	s.command = c
	if s.setupErr == nil {
		s.state.IsSetup = true
	}
	return s.setupErr
}

type setupPageConfig struct {
	siteconfig.Store
	data siteconfig.ConfigData
}

func (s *setupPageConfig) Get() *siteconfig.ConfigData { return &s.data }

func setupPageRouter(s *setupPageService) http.Handler {
	log := zap.NewNop().Sugar()
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../translations/en.toml"}}, log)
	r := chi.NewRouter()
	r.Use(olhttp.ReqCtxMiddleware, provider.Middleware, csrf.NewHandler("test-secret").Middleware)
	frontend.AttachAssetsInliningHandler(os.DirFS("../frontend/embed-assets"), "embed-assets", r)
	config := &setupPageConfig{}
	config.data.PasswordRequirements.MinLength = 8
	middleware := newSetupMiddleware(s, config, log)
	r.Use(middleware.Middleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	return r
}

func setupPageGet(h http.Handler, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
	return w
}
func setupPagePost(h http.Handler, values url.Values, validCSRF bool) *httptest.ResponseRecorder {
	get := setupPageGet(h, "/setup")
	req := httptest.NewRequest("POST", "/setup", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if validCSRF {
		for _, cookie := range get.Result().Cookies() {
			req.AddCookie(cookie)
			if cookie.Name == "csrf" {
				req.Header.Set("x-csrf-token", cookie.Value)
			}
		}
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestSetupPageForms(t *testing.T) {
	for _, admin := range []bool{false, true} {
		s := &setupPageService{state: app.GetSetupDataResult{AdminUserExist: admin}}
		w := setupPageGet(setupPageRouter(s), "/setup")
		body := w.Body.String()
		if w.Code != 200 || !strings.Contains(body, `name="__csrf"`) || !strings.Contains(body, `name="siteName"`) {
			t.Fatalf("bad form: %d %s", w.Code, body)
		}
		if strings.Contains(body, `name="adminPassword"`) == admin {
			t.Fatal("admin creation visibility incorrect")
		}
	}
}
func TestSetupPageSubmission(t *testing.T) {
	for _, admin := range []bool{false, true} {
		s := &setupPageService{state: app.GetSetupDataResult{AdminUserExist: admin}}
		h := setupPageRouter(s)
		values := url.Values{"siteName": {" My Library "}}
		if !admin {
			values.Set("adminUsername", "owner")
			values.Set("adminPassword", "password123")
			values.Set("confirmPassword", "password123")
		}
		w := setupPagePost(h, values, true)
		if w.Code != 303 || w.Header().Get("Location") != "/login" || s.calls != 1 || s.command.SiteName != "My Library" {
			t.Fatalf("bad submission: %d %+v", w.Code, s)
		}
		if w := setupPageGet(h, "/"); w.Code != 204 {
			t.Fatal("setup still intercepts completed site", w.Code)
		}
		if w := setupPageGet(h, "/setup"); w.Code != 303 {
			t.Fatal("completed setup remains accessible", w.Code)
		}
	}
}
func TestSetupPageValidationAndCSRF(t *testing.T) {
	s := &setupPageService{}
	h := setupPageRouter(s)
	values := url.Values{"siteName": {"Library"}, "adminUsername": {"owner"}, "adminPassword": {"do-not-echo"}, "confirmPassword": {"different"}}
	w := setupPagePost(h, values, true)
	if w.Code != 422 || s.calls != 0 || !strings.Contains(w.Body.String(), "passwords do not match") || strings.Contains(w.Body.String(), "do-not-echo") || !strings.Contains(w.Body.String(), `value="owner"`) {
		t.Fatalf("bad validation: %d %s", w.Code, w.Body.String())
	}
	w = setupPagePost(h, values, false)
	if w.Code != 403 || s.calls != 0 {
		t.Fatal("csrf bypass", w.Code, s.calls)
	}
}
func TestSetupPageFailures(t *testing.T) {
	s := &setupPageService{checkErr: errors.New("private database details")}
	h := setupPageRouter(s)
	for _, path := range []string{"/", "/setup"} {
		w := setupPageGet(h, path)
		if w.Code != 503 || strings.Contains(w.Body.String(), "private database") {
			t.Fatal("state failure leaked or allowed traffic", w.Code)
		}
	}
	s.checkErr = nil
	s.state.AdminUserExist = true
	s.setupErr = errors.New("private database details")
	w := setupPagePost(h, url.Values{"siteName": {"Library"}}, true)
	if w.Code != 500 || strings.Contains(w.Body.String(), "private database") || !strings.Contains(w.Body.String(), "Please try again") {
		t.Fatal("save failure mishandled", w.Code, w.Body.String())
	}
	s.setupErr = app.ErrUsernameTaken
	w = setupPagePost(h, url.Values{"siteName": {"Library"}}, true)
	if w.Code != 422 || !strings.Contains(w.Body.String(), "already taken") {
		t.Fatal("username conflict mishandled", w.Code)
	}
}
func TestSetupRedirectAndMethods(t *testing.T) {
	h := setupPageRouter(&setupPageService{})
	if w := setupPageGet(h, "/"); w.Code != 303 || w.Header().Get("Location") != "/setup" {
		t.Fatal("missing setup redirect", w.Code)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("OPTIONS", "/setup", nil))
	if w.Code != 405 || w.Header().Get("Allow") != "GET, POST" {
		t.Fatal("wrong methods", w.Code)
	}
}
