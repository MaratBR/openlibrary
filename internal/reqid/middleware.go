package reqid

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type keyType struct{}

var key keyType

var instanceID string

func init() {
	ts := uint64(time.Now().Unix())
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], ts)
	instanceID = base64.RawURLEncoding.EncodeToString(b[:])
}

func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := generateID()
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

var globalRequestID int64

func generateID() string {
	intID := atomic.AddInt64(&globalRequestID, 1)
	return instanceID + strconv.FormatInt(intID, 36)
}
