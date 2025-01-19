package publicui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/gofrs/uuid"
)

func write500(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func writeRequestError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}

func writeApplicationError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(409)
	w.Write([]byte(err.Error()))
}

func writeUnauthorizedError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}

func readJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func getInt64Array(value url.Values, key string) []int64 {
	strArr := getStringArray(value, key)

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

func getUUIDArray(value url.Values, key string) []uuid.UUID {
	strArr := getStringArray(value, key)

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

func getInt32FromQuery(values url.Values, key string) app.Int32 {
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

func getInt32RangeFromQuery(query url.Values, queryParam string) app.Int32Range {
	max := getInt32FromQuery(query, fmt.Sprintf("%s.max", queryParam))
	min := getInt32FromQuery(query, fmt.Sprintf("%s.min", queryParam))
	return app.Int32Range{Max: max, Min: min}
}

func getPage(values url.Values, key string) uint {
	v := getInt32FromQuery(values, key)

	if !v.Valid {
		return 1
	}

	if v.Int32 < 1 {
		return 1
	}

	return uint(v.Int32)
}

func getStringArray(value url.Values, key string) []string {
	str := value.Get(key)

	if str == "" {
		return nil
	}

	return splitByWithEscape(str, '|')
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
