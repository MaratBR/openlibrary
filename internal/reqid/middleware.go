package reqid

import (
	"net/http"

	"github.com/gofrs/uuid"
)

func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := generateID(r)

			w.Header().Add("X-ReqID", id)

			next.ServeHTTP(w, r)
		})
	}
}

func generateID(r *http.Request) string {
	return uuid.Must(uuid.NewV4()).String()
}
