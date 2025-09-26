package public

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type authController struct {
	authService app.AuthService
	csrfHandler *csrf.Handler
}

func newAuthController(authService app.AuthService, csrfHandler *csrf.Handler) *authController {
	return &authController{authService: authService, csrfHandler: csrfHandler}
}

func (c *authController) Register(r chi.Router) {
	// public auth pages
	r.HandleFunc("/login", c.LogIn)
	r.HandleFunc("/logout", c.LogOut)
}

func (c *authController) LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, ok := auth.GetSession(r.Context())
		if ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		c.writeLoginForm(w, r, templates.LoginData{})
	} else if r.Method == http.MethodPost {
		c.handleSignIn(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	templates.Login(data).Render(r.Context(), w)
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

func (c *authController) LogOut(w http.ResponseWriter, r *http.Request) {
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
