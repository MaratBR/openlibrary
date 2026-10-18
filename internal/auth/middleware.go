package auth

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/internal/store"
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
	log *zap.SugaredLogger,
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
			sessionInstance, ok := session.Get(r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			// get session info
			sessionID := sessionInstance.ID()

			// Session revocation and moderation bans must take effect on the next request.
			sessionInfo, err := sessionService.GetBySID(r.Context(), sessionID)
			if err != nil {
				if err == app.ErrSessionNotFound || err == app.ErrUserBanned {
					session.ResetSession(r, w)
					next.ServeHTTP(w, r)
				} else {
					log.Errorw("failed to retrieve session", "err", err)
					options.OnFail(w, r, err)
				}
				return
			}

			user, err := userService.GetUserSelfData(r.Context(), sessionInfo.UserID)
			if err != nil {
				if err == store.ErrNoRows {
					// user with this ID does not exist and session info was likely
					// retrieved from expired cache
					log.Warnw("detected expired cache data, resetting user session", "userID", sessionInfo.UserID)
					session.ResetSession(r, w)
					next.ServeHTTP(w, r)
					return
				}
				log.Errorw("unexpected error when trying to retrieve user's data", "err", err)
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
