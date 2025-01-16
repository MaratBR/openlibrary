package httputil

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func URLParamInt64(r *http.Request, name string) (int64, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func URLParamUUID(r *http.Request, name string) (uuid.UUID, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		return uuid.Nil, nil
	}
	id, err := uuid.FromString(value)
	return id, err
}
