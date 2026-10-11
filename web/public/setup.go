package public

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/siteconfig"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"go.uber.org/zap"
)

type setupMiddleware struct {
	setupService app.SetupService
	siteConfig   siteconfig.Store
	log          *zap.SugaredLogger
}

func newSetupMiddleware(service app.SetupService, config siteconfig.Store, log *zap.SugaredLogger) *setupMiddleware {
	return &setupMiddleware{setupService: service, siteConfig: config, log: log}
}

func (s *setupMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/setup" {
			w.Header().Set("Cache-Control", "no-store")
			s.intercept(w, r)
			return
		}
		done, err := s.setupService.IsSetup()
		if err != nil {
			s.log.Errorw("failed to check server setup state", "err", err)
			http.Error(w, i18n.GetLocalizer(r.Context()).T("setup.unavailable"), http.StatusServiceUnavailable)
			return
		}
		if !done {
			w.Header().Set("Cache-Control", "no-store")
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *setupMiddleware) intercept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		methodNotAllowed(w, r)
		return
	}
	state, err := s.setupService.GetSetupState()
	if err != nil {
		s.log.Errorw("failed to load setup state", "err", err)
		http.Error(w, i18n.GetLocalizer(r.Context()).T("setup.unavailable"), http.StatusServiceUnavailable)
		return
	}
	if state.IsSetup {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	data := templates.InitialSetupData{AdminUserExist: state.AdminUserExist, SiteName: s.siteConfig.Get().SiteSetup.SiteName, PasswordRequirements: s.siteConfig.Get().PasswordRequirements}
	status := http.StatusOK
	if r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		if err := r.ParseForm(); err != nil {
			writeBadRequest(w, r, err)
			return
		}
		data.SiteName = strings.TrimSpace(r.PostForm.Get("siteName"))
		data.AdminUsername = r.PostForm.Get("adminUsername")
		password := r.PostForm.Get("adminPassword")
		switch {
		case data.SiteName == "":
			data.ErrorKey = "setup.siteNameRequired"
		case !state.AdminUserExist && password != r.PostForm.Get("confirmPassword"):
			data.ErrorKey = "setup.passwordMismatch"
		default:
			err = s.setupService.Setup(r.Context(), app.SetupCommand{SiteName: data.SiteName, AdminUsername: data.AdminUsername, AdminPassword: password})
			if err == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			switch {
			case errors.Is(err, app.ErrSetupInvalidInput):
				data.ErrorKey = "setup.invalidInput"
			case errors.Is(err, app.ErrUsernameTaken):
				data.ErrorKey = "setup.usernameTaken"
			default:
				s.log.Errorw("failed to complete setup", "err", err)
				data.ErrorKey = "setup.failed"
				status = http.StatusInternalServerError
			}
			// Admin creation may have succeeded before configuration persistence failed.
			if fresh, readErr := s.setupService.GetSetupState(); readErr == nil {
				data.AdminUserExist = fresh.AdminUserExist
			}
		}
		if status == http.StatusOK {
			status = http.StatusUnprocessableEntity
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	olhttp.WriteTemplate(w, r.Context(), templates.InitialSetup(data))
}
