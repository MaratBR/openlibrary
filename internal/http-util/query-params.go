package httputil

import (
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
)

func URLQueryParamInt64(r *http.Request, name string) (int64, error) {
	value := r.URL.Query().Get(name)
	if len(value) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func URLQueryParamUUID(r *http.Request, name string) (uuid.UUID, error) {
	value := r.URL.Query().Get(name)
	if len(value) == 0 {
		return uuid.Nil, nil
	}
	return uuid.FromString(value)
}
