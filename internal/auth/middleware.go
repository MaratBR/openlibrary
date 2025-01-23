package auth

import (
	"log/slog"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
)

type MiddlewareOptions struct {
	When   func(r *http.Request) bool
	OnFail func(w http.ResponseWriter, r *http.Request, err error)
}

func NewAuthorizationMiddleware(service app.SessionService, options MiddlewareOptions) func(http.Handler) http.Handler {
	if options.OnFail == nil {
		options.OnFail = defaultAuthorizationFailedHandler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.When != nil && !options.When(r) {
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
					options.OnFail(w, r, err)
				}
				return
			}

			r = AttachSessionInfo(r, sessionInfo)
			next.ServeHTTP(w, r)
		})
	}
}

func defaultAuthorizationFailedHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte("Authorization failed\n"))
	w.Write([]byte(err.Error()))
}
