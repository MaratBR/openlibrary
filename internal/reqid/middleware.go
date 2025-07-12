package reqid

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
)

type keyType struct{}

var key keyType

func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := generateID(r)

			w.Header().Add("X-ReqID", id)
			r = r.WithContext(context.WithValue(r.Context(), key, id))

			next.ServeHTTP(w, r)
		})
	}
}

func Get(r *http.Request) string {
	v := r.Context().Value(key)
	if v == nil {
		return ""
	}
	return v.(string)
}

func generateID(r *http.Request) string {
	return uuid.Must(uuid.NewV4()).String()
}
