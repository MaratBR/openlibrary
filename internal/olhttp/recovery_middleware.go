package olhttp

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
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

					stack := debug.Stack()

					fmt.Println("--- Custom Panic Handler ---")
					fmt.Printf("Recovered from panic: %v\n", r)
					fmt.Println("Stack trace:")
					os.Stderr.Write(stack) // Print the stack trace
					fmt.Println("--------------------------")

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
