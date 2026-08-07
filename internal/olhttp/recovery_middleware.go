package olhttp

import (
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/MaratBR/openlibrary/internal/reqid"
)

func MakeRecoveryMiddleware(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// TODO add request id
					stack := debug.Stack()
					log.Errorw("recovered from panic", "panic", rec, "requestID", reqid.Get(r), "stack", string(stack))

					b := strings.Builder{}
					b.WriteString("server panicked!\nIf you are a developer, please fix this. If not please contact support or report to https://github.com/MaratBR/openlibrary/issues/new\n")
					b.WriteString("\n\nrequest id: ")
					b.WriteString(reqid.Get(r))
					b.WriteString("\n\nSTACK TRACE:\n")
					b.Write(stack)
					body := b.String()
					w.Write([]byte(body))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
