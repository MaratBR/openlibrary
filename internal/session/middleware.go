package session

import (
	"context"
	"log/slog"
	"net/http"
)

type ctxKey uint8

const (
	ctxKeySession ctxKey = 1
)

const CookieName = "sid"

func Middleware(
	store Store,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionIDCookie, err := r.Cookie(CookieName)
			if err != nil {
				if err != http.ErrNoCookie {
					// some kind of weird error
					slog.Error("unknown error while trying to get session id", "cookie", CookieName)
				}

				// no cookie - just continue
				next.ServeHTTP(w, r)
				return
			}

			sessionID := sessionIDCookie.Value

			if sessionID == "" {
				// empty session ID - just ignore that
				next.ServeHTTP(w, r)
				return
			}

			session, err := store.Get(r.Context(), sessionID)

			if err != nil {
				if err != ErrNoSession {
					slog.Error("unknown error while trying to get a session info", "err", err)
				}

				next.ServeHTTP(w, r)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), ctxKeySession, session))
			next.ServeHTTP(w, r)

			err = session.Save(r.Context())
			if err != nil {
				slog.Error("failed to save session", "err", err)
			}
		})
	}
}

func Get(r *http.Request) (Session, bool) {
	v := r.Context().Value(ctxKeySession)
	if v == nil {
		return nil, false
	}

	return v.(Session), true
}

func ResetSession(r *http.Request, w http.ResponseWriter) {
	RemoveSIDCookie(w, isSecure(r))
	// TODO: when session is reset - replace it with a new one for anonymous user
}

func isSecure(r *http.Request) bool {
	return r.URL.Scheme == "https"
}
