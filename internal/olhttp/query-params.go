package olhttp

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app"
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

func GetInt32FromQuery(values url.Values, key string) app.Int32 {
	v := values.Get(key)
	if v == "" {
		return app.Int32{}
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return app.Int32{}
	}
	return app.Int32{Valid: true, Int32: int32(i)}
}

func GetPage(values url.Values, key string) uint32 {
	v := GetInt32FromQuery(values, key)

	if !v.Valid {
		return 1
	}

	if v.Int32 < 1 {
		return 1
	}

	return uint32(v.Int32)
}

func GetPageSize(values url.Values, key string, minValue, maxValue, defaultValue uint32) uint32 {
	v := GetInt32FromQuery(values, key)

	if !v.Valid {
		return defaultValue
	}

	if v.Int32 < int32(minValue) {
		return minValue
	}

	if v.Int32 > int32(maxValue) {
		return maxValue
	}

	return uint32(v.Int32)
}

func GetBool(value url.Values, key string) app.Nullable[bool] {
	v := value.Get(key)
	if v == "" {
		return app.Null[bool]()
	}

	if v == "true" || v == "on" || v == "1" || v == "yes" || v == "y" {
		return app.Value(true)
	}

	if v == "false" || v == "off" || v == "0" || v == "no" || v == "n" {
		return app.Value(false)
	}

	return app.Null[bool]()
}

func GetBoolDefault(value url.Values, key string, defaultValue bool) bool {
	v := GetBool(value, key)
	if v.Valid {
		return v.Value
	} else {
		return defaultValue
	}
}
