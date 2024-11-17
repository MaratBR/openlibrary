package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/store"
)

type authController struct {
	authService *app.AuthService
}

func newAuthController(service *app.AuthService) authController {
	return authController{authService: service}
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c authController) SignIn(w http.ResponseWriter, r *http.Request) {
	req, err := getJSON[signInRequest](r)
	if err != nil {
		writeRequestError(err, w)
		return
	}
	if req.Username == "" || req.Password == "" || len(req.Username) > 50 || len(req.Password) > 50 {
		writeUnprocessableEntity(w, "username and password both must be between 1 and 50 characters")
		return
	}
	c.signIn(req.Username, req.Password, w, r)
}

func (c authController) signIn(
	username, password string,
	w http.ResponseWriter, r *http.Request) {
	result, err := c.authService.SignIn(r.Context(), app.SignInCommand{
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	httpSecure := "; Secure"
	if r.URL.Scheme != "https" {
		httpSecure = ""
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("sid=%s; Path=/; Max-Age=%d; HttpOnly%s", result.SessionID, int((time.Hour*24*365).Seconds()), httpSecure))
	w.Write([]byte("OK"))
}

type signUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c authController) SignUp(w http.ResponseWriter, r *http.Request) {
	c.authService.SignIn(r.Context(), app.SignInCommand{
		Username: "admin",
		Password: "admin",
	})
}

func (c authController) SignInAdmin(w http.ResponseWriter, r *http.Request) {
	c.signIn("admin", "admin", w, r)
}

func refreshCsrfToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func newAuthorizationMiddlewareConditional(service *app.AuthService, when func(r *http.Request) bool) func(http.Handler) http.Handler {
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

			sessionInfo, err := service.GetUserBySessionID(r.Context(), sidCookie.Value)
			if err != nil {
				if err == store.ErrNoRows {
					next.ServeHTTP(w, r)
				} else {
					slog.Error("unexpected error when trying to retrieve user's session", "err", err)
					next.ServeHTTP(w, r)
				}
				return
			}

			newContext := context.WithValue(r.Context(), sessionInfoKey, sessionInfo)
			r = r.WithContext(newContext)
			w.Header().Add("x-user-id", sessionInfo.UserID.String())
			next.ServeHTTP(w, r)
		})
	}
}

func newAuthorizationMiddleware(service *app.AuthService) func(http.Handler) http.Handler {
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
