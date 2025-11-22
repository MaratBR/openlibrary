package olhttp

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

var (
	errNsUrlParams             = httpErrors.NewSubNamespace("url_params")
	errTypeInvalidInt64        = errNsUrlParams.NewType("int64")
	errInvalidInt64EmptyString = errTypeInvalidInt64.New("invalid int64: empty string")
	errTypeInvalidUUID         = errNsUrlParams.NewType("uuid")
)

func URLParamInt64(r *http.Request, name string) (int64, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		return 0, errInvalidInt64EmptyString
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errTypeInvalidInt64.Wrap(err, "failed to parse int64 url parameter")
	}
	return intValue, nil
}

func URLParamUUID(r *http.Request, name string) (uuid.UUID, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		return uuid.Nil, nil
	}
	id, err := uuid.FromString(value)
	if err != nil {
		return uuid.Nil, errTypeInvalidUUID.Wrap(err, "failed to parse uuid url parameter")
	}
	return id, nil
}

func IsSameURL(url *url.URL, str string) bool {
	if str == "" {
		return false
	}

	if strings.HasPrefix(str, "?") {
		return "?"+url.Query().Encode() == str
	}

	return url.String() == str
}
