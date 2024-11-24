package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/cmd/server/csrf"
	"github.com/MaratBR/openlibrary/internal/app"
)

type authController struct {
	authService    app.AuthService
	bypassTLSCheck bool
	csrfHandler    *csrf.Handler
}

func (c authController) requireTLS(r *http.Request, w http.ResponseWriter) bool {
	if c.bypassTLSCheck {
		return true
	}

	if r.TLS == nil {
		writeTLSRequiredError(w)
		return false
	}

	return true
}

func newAuthController(service app.AuthService, bypassTLSCheck bool, csrfHandler *csrf.Handler) authController {
	return authController{authService: service, bypassTLSCheck: bypassTLSCheck, csrfHandler: csrfHandler}
}

func (c authController) SignOut(w http.ResponseWriter, r *http.Request) {
	if !c.requireTLS(r, w) {
		return
	}

	session, ok := getSession(r)

	if !ok {
		writeOK(w)
		return
	}

	err := c.authService.SignOut(r.Context(), session.SessionID)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	c.csrfHandler.WriteAnonymousCSRFToken(w)
	removeSidCookie(w, "sid", r.URL.Scheme == "https")

	writeOK(w)
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c authController) SignIn(w http.ResponseWriter, r *http.Request) {
	if !c.requireTLS(r, w) {
		return
	}

	req, err := getJSON[signInRequest](r)
	if err != nil {
		writeRequestError(err, w)
		return
	}
	if len(req.Username) == 0 || len(req.Username) > 50 {
		writeUnprocessableEntity(w, "username must be between 1 and 50 characters")
		return
	}
	if len(req.Password) == 0 || len(req.Password) > 200 {
		writeUnprocessableEntity(w, "password must be between 1 and 200 characters")
		return
	}
	c.signIn(req.Username, req.Password, w, r)
}

func (c authController) signIn(
	username, password string,
	w http.ResponseWriter, r *http.Request) {
	result, err := c.authService.SignIn(r.Context(), app.SignInCommand{
		Username: username,
		Password: password,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeSidCookie(w, c.csrfHandler.SIDCookie(), result.SessionID, time.Hour*24*30, r.URL.Scheme == "https")
	c.csrfHandler.WriteCSRFToken(w, result.SessionID)
	w.Write([]byte("OK"))
}

type signUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c authController) SignUp(w http.ResponseWriter, r *http.Request) {
	if !c.requireTLS(r, w) {
		return
	}

	req, err := getJSON[signUpRequest](r)
	if err != nil {
		writeRequestError(err, w)
		return
	}

	if len(req.Username) == 0 || len(req.Username) > 50 {
		writeUnprocessableEntity(w, "username must be between 1 and 50 characters")
		return
	}
	if len(req.Password) == 0 || len(req.Password) > 200 {
		writeUnprocessableEntity(w, "password must be between 1 and 200 characters")
		return
	}

	result, err := c.authService.SignUp(r.Context(), app.SignUpCommand{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: r.UserAgent(),
		IpAddress: r.RemoteAddr,
	})

	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeSidCookie(w, "sid", result.SessionID, time.Hour*24*30, r.URL.Scheme == "https")
	c.csrfHandler.WriteCSRFToken(w, result.SessionID)

	writeOK(w)
}

func (c authController) SignInAdmin(w http.ResponseWriter, r *http.Request) {
	if !c.requireTLS(r, w) {
		return
	}

	c.signIn("admin", "admin", w, r)
}

func newAuthorizationMiddlewareConditional(service app.SessionService, when func(r *http.Request) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !when(r) {
				next.ServeHTTP(w, r)
				return
			}

			sidCookie, err := r.Cookie("sid")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			sessionInfo, err := service.GetBySID(r.Context(), sidCookie.Value)
			if err != nil {
				if err == app.ErrSessionNotFound {
					next.ServeHTTP(w, r)
				} else {
					slog.Error("unexpected error when trying to retrieve user's session", "err", err)
					writeApplicationError(w, err)
				}
				return
			}

			newContext := context.WithValue(r.Context(), sessionInfoKey, sessionInfo)
			r = r.WithContext(newContext)
			next.ServeHTTP(w, r)
		})
	}
}

func newAuthorizationMiddleware(service app.SessionService) func(http.Handler) http.Handler {
	return newAuthorizationMiddlewareConditional(service, func(r *http.Request) bool { return true })
}

func newRequireAuthorizationMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := getSession(r)
			if !ok {
				writeUnauthorizedError(w)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

type sessionInfoKeyType struct{}

var sessionInfoKey sessionInfoKeyType = sessionInfoKeyType{}
