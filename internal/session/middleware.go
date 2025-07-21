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

func Middleware(
	cookieName string,
	store Store,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionIDCookie, err := r.Cookie(cookieName)
			if err != nil {
				if err != http.ErrNoCookie {
					// some kind of weird error
					slog.Error("unknown error while trying to get session id", "cookie", cookieName)
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
