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

func splitByWithEscape(s string, c byte) []string {
	result := []string{}
	buf := []byte{}
	escaped := false

	for i := 0; i < len(s); i++ {
		if escaped {
			escaped = false
			buf = append(buf, s[i])
			continue
		}

		if s[i] == '\\' {
			escaped = true
			continue
		}

		if s[i] == c {
			result = append(result, string(buf))
			buf = nil
			continue
		}

		buf = append(buf, s[i])
	}

	if len(buf) > 0 {
		result = append(result, string(buf))
	}

	return result
}

func GetStringArray(value url.Values, key string) []string {
	values, ok := value[key]

	if !ok || len(values) == 0 {
		return nil
	}

	if len(values) == 1 {
		return splitByWithEscape(values[0], ',')
	} else {
		return values
	}
}

func GetInt64Array(value url.Values, key string) []int64 {
	strArr := GetStringArray(value, key)

	if strArr == nil {
		return nil
	}

	i64Arr := []int64{}
	for _, str := range strArr {
		id, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			continue
		}

		i64Arr = append(i64Arr, id)
	}

	return i64Arr

}

func GetUUIDArray(value url.Values, key string) []uuid.UUID {
	strArr := GetStringArray(value, key)

	if strArr == nil {
		return nil
	}

	uuidArr := []uuid.UUID{}
	for _, str := range strArr {
		id, err := uuid.FromString(str)
		if err != nil {
			continue
		}

		uuidArr = append(uuidArr, id)
	}

	return uuidArr
}
