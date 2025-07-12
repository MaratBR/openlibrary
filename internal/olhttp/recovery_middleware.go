package olhttp

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/MaratBR/openlibrary/internal/reqid"
)

func MakeRecoveryMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// TODO add request id
					slog.ErrorContext(r.Context(), "recovered from panic", "rec", rec, "request-id", reqid.Get(r))

					b := strings.Builder{}
					b.WriteString("server panicked!\nIf you are a developer, please fix this. If not please contact support or report to https://github.com/MaratBR/openlibrary/issues/new\n")
					b.WriteString("\n\nrequest id: ")
					b.WriteString(reqid.Get(r))
					body := b.String()
					w.Write([]byte(body))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
