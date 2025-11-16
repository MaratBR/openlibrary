package public

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
)

type authController struct {
	authService   app.AuthService
	signUpService app.SignUpService
	csrfHandler   *csrf.Handler
	siteConfig    *app.SiteConfig
	cfg           *koanf.Koanf
}

func newAuthController(authService app.AuthService, signUpService app.SignUpService, csrfHandler *csrf.Handler, siteConfig *app.SiteConfig, cfg *koanf.Koanf) *authController {
	return &authController{authService: authService, csrfHandler: csrfHandler, siteConfig: siteConfig, cfg: cfg, signUpService: signUpService}
}

func (c *authController) Register(r chi.Router) {
	// public auth pages
	r.HandleFunc("/login", c.login)
	r.HandleFunc("/logout", c.logout)
	r.HandleFunc("/signup", c.signup)
	r.HandleFunc("/signup/email-verification-code", c.emailVerification)
}

func (c *authController) emailVerification(w http.ResponseWriter, r *http.Request) {
	user, isAuthorized := auth.GetUser(r.Context())
	switch r.Method {
	case http.MethodGet:
		if !isAuthorized || user.IsEmailVerified {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		status, err := c.signUpService.GetEmailVerificationStatus(r.Context(), user.ID)
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}
		olhttp.WriteTemplate(w, r.Context(), templates.SignUpEmailVerification(user.Email, status))
	case http.MethodPost:
		if !isAuthorized {
			writeUnauthorizedError(w)
			return
		}

		err := c.signUpService.VerifyEmail(r.Context(), app.VerifyEmailCommand{
			Code:   r.FormValue("code"),
			UserID: user.ID,
		})
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	case http.MethodPatch:
		act := r.URL.Query().Get("act")
		if act == "resend" {
			result, err := c.signUpService.SendEmailVerification(r.Context(), app.SendEmailVerificationCommand{
				UserID:          user.ID,
				BypassRateLimit: false,
			})
			if err != nil {
				apiWriteApplicationError(w, err)
				return
			}
			var resp struct {
				CanResendAfter time.Time `json:"canResendAfter"`
			}
			resp.CanResendAfter = result.CanResendAfter
			olhttp.NewAPIResponse(resp).Write(w)
		} else {
			apiWriteUnprocessableEntity(w, errors.New("unknown act value"))
		}
	default:
		methodNotAllowed(w, r)
	}
}

func (c *authController) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, ok := auth.GetSession(r.Context())
		if ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		c.writeLoginForm(w, r, templates.LoginData{})
	case http.MethodPost:
		c.handleSignIn(w, r)
	default:
		methodNotAllowed(w, r)
	}
}

func (c *authController) writeLoginForm(w http.ResponseWriter, r *http.Request, data templates.LoginData) {
	if data.InitialLogin == "" {
		initialLoginValue, err := r.Cookie("auth_ll")
		if err == nil {
			data.InitialLogin = initialLoginValue.Value
		}
	}

	if r.URL.Query().Get("next") == "/admin" {
		data.IsToAdmin = true
	}

	olhttp.WriteTemplate(w, r.Context(), templates.Login(data))
}

func (c *authController) handleSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 {
		c.writeLoginForm(w, r, templates.LoginData{Error: "username must not be empty"})
		return
	}

	c.signIn(username, password, w, r)
}

func (c *authController) signIn(username string, password string, w http.ResponseWriter, r *http.Request) {
	result, err := c.authService.SignIn(r.Context(), app.SignInCommand{
		Username:  username,
		Password:  password,
		UserAgent: r.UserAgent(),
		IpAddress: r.RemoteAddr,
	})
	if err != nil {
		c.writeLoginForm(w, r, templates.LoginData{Error: err.Error()})
		return
	}
	session.WriteSIDCookie(w, result.SessionID, time.Hour*24*30, r.URL.Scheme == "https")
	c.csrfHandler.WriteCSRFToken(w, result.SessionID)
	w.Header().Add("Set-Cookie", "auth_ll="+username)

	c.redirectToNext(w, r)
}

func (c *authController) redirectToNext(w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("next")
	if next == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	u, err := url.Parse(next)
	if err != nil {
		slog.Warn("failed to parse next param", "err", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if u.Scheme != "" || u.Host != "" || u.User != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, next, http.StatusFound)
	}
}

func (c *authController) logout(w http.ResponseWriter, r *http.Request) {
	session, ok := auth.GetSession(r.Context())
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err := c.authService.SignOut(r.Context(), session.SID)
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	w.Header().Add("Set-Cookie", "sid=deleted; expires=Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Add("Set-Cookie", "auth_ll=deleted; expires=Thu, 01 Jan 1970 00:00:00 GMT")
	c.csrfHandler.WriteAnonymousCSRFToken(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *authController) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, isAuthorized := auth.GetSession(r.Context())
		if isAuthorized {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		olhttp.WriteTemplate(w, r.Context(), templates.SignUp(c.siteConfig, c.cfg.Bool("auth.requireEmail")))
	case http.MethodPost:
		c.handleSignUp(w, r)
	default:
		methodNotAllowed(w, r)
	}
}

func (c *authController) handleSignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	userName := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	ip := olhttp.GetIP(r).String()
	userAgent := r.UserAgent()
	result, err := c.signUpService.SignUp(r.Context(), app.SignUpCommand{
		Username:  userName,
		Password:  password,
		Email:     email,
		IpAddress: ip,
		UserAgent: userAgent,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	// create session and send them on their way
	if result.Created {
		sessionID, err := c.authService.CreateSessionForUser(r.Context(), result.CreatedUserID, userAgent, ip)
		if err != nil {
			// TODO redirect to login page and let them login?
			writeApplicationError(w, r, err)
			return
		}

		session.WriteSIDCookie(w, sessionID, time.Hour*24*30, r.URL.Scheme == "https")
		if result.EmailVerificationRequired {
			http.Redirect(w, r, "/signup/email-verification-code", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		if result.EmailTaken {
			// TODO impelement a email reclaim page or something
			write500(w, r, errors.New("dev is lazy and did not implement this yet"))
			return
		}

		// some weird stuff is going on
		write500(w, r, errors.New("unknown error while trying to create a user, if you see this that means some really weird stuff is going on and the developer of this site is probably the one to blaim (but please be nice ok?)"))
	}

}
