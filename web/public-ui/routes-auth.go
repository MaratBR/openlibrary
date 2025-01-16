package publicui

import (
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/web/public-ui/templates"
)

type authController struct {
	authService app.AuthService
	csrfHandler *csrf.Handler
}

func newAuthController(authService app.AuthService, csrfHandler *csrf.Handler) *authController {
	return &authController{authService: authService, csrfHandler: csrfHandler}
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

	templates.Login(r.Context(), data).Render(r.Context(), w)
}

func (c *authController) handleSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeRequestError(w, r, err)
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
	auth.WriteSIDCookie(w, c.csrfHandler.SIDCookie(), result.SessionID, time.Hour*24*30, r.URL.Scheme == "https")
	c.csrfHandler.WriteCSRFToken(w, result.SessionID)
	w.Header().Add("Set-Cookie", "auth_ll="+username)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *authController) LogOut(w http.ResponseWriter, r *http.Request) {}
