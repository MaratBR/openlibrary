package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/session"
)

var (
	_CTX_SESSION_INFO_SESSION_KEY = "ol-session"
)

type MiddlewareOptions struct {
	When   func(r *http.Request) bool
	OnFail func(w http.ResponseWriter, r *http.Request, err error)
}

func NewAuthorizationMiddleware(
	sessionService app.SessionService,
	userService app.UserService,
	options MiddlewareOptions,
) func(http.Handler) http.Handler {
	if options.OnFail == nil {
		options.OnFail = defaultAuthorizationFailedHandler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// filter out auth from some endpoints
			if options.When != nil && !options.When(r) {
				next.ServeHTTP(w, r)
				return
			}

			// get session instance
			sessionInstance, ok := session.GetSession(r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			// get session info
			sessionID := sessionInstance.ID()

			var sessionInfo *app.SessionInfo

			{
				stringValue, ok := sessionInstance.Get(_CTX_SESSION_INFO_SESSION_KEY)

				var err error
				if !ok {
					sessionInfo, err = sessionService.GetBySID(r.Context(), sessionID)
					if err != nil {
						if err == app.ErrSessionNotFound {
							next.ServeHTTP(w, r)
						} else {
							slog.Error("unexpected error when trying to retrieve user's session", "err", err)
							options.OnFail(w, r, err)
						}
						return
					}

					b, err := json.Marshal(sessionInfo)
					if err == nil {
						sessionInstance.Put(_CTX_SESSION_INFO_SESSION_KEY, string(b))
						err = sessionInstance.Save(r.Context())
						if err != nil {
							slog.Error("error while saving session", "err", err)
						}
					}
				} else {
					sessionInfo = new(app.SessionInfo)
					err = json.Unmarshal([]byte(stringValue), sessionInfo)
					if err != nil {
						// perhaps session was corrupted or the schema changed
						slog.Error("error while parsing session info", "err", err)
						next.ServeHTTP(w, r)
					}
				}
			}

			user, err := userService.GetUserSelfData(r.Context(), sessionInfo.UserID)
			if err != nil {
				slog.Error("unexpected error when trying to retrieve user's data", "err", err)
				options.OnFail(w, r, err)
				return
			}

			r = attachSessionInfo(r, sessionInfo, user)
			next.ServeHTTP(w, r)
		})
	}
}

func defaultAuthorizationFailedHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte("Authorization failed\n"))
	w.Write([]byte(err.Error()))
}
