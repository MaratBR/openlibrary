package olhttp

import (
	"log/slog"
	"net/http"
)

func MakeRecoveryMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// TODO add request id
					slog.ErrorContext(r.Context(), "recovered from panic", "rec", rec)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
