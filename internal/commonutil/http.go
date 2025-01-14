package commonutil

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func URLParamInt64(r *http.Request, name string) (int64, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}
