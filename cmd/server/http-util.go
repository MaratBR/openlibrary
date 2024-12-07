package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func readJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func getJSON[T any](r *http.Request) (T, error) {
	var value T
	err := readJSON(r, &value)
	if err != nil {
		return value, jsonBodyError{err}
	}
	return value, nil
}

type jsonBodyError struct {
	err error
}

func (err jsonBodyError) Error() string {
	return err.Error()
}

func urlParamInt64(r *http.Request, name string) (int64, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func urlParamUUID(r *http.Request, name string) (uuid.UUID, error) {
	value := chi.URLParam(r, name)
	if len(value) == 0 {
		return uuid.Nil, nil
	}
	id, err := uuid.FromString(value)
	return id, err
}

func urlQueryParamInt64(r *http.Request, name string) (int64, error) {
	value := r.URL.Query().Get(name)
	if len(value) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func urlQueryParamUUID(r *http.Request, name string) (uuid.UUID, error) {
	value := r.URL.Query().Get(name)
	if len(value) == 0 {
		return uuid.Nil, nil
	}
	return uuid.FromString(value)
}

func unwrapInt64StringArr(arr []app.Int64String) []int64 {
	arr2 := make([]int64, len(arr))
	for i, v := range arr {
		arr2[i] = int64(v)
	}
	return arr2
}

func readUrlEncodedBody(r *http.Request) (url.Values, error) {
	if r.ContentLength > 20_000 {
		return nil, ErrHttpRequestBodyTooLarge
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	return values, nil
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

func stringArray(arr []string) string {
	if len(arr) == 0 {
		return ""
	}

	arr2 := make([]string, len(arr))
	for i, v := range arr {
		arr2[i] = strings.ReplaceAll(v, "|", "\\|")
	}

	return strings.Join(arr2, "|")
}

func i64Array(arr []int64) string {
	if len(arr) == 0 {
		return ""
	}

	arr2 := make([]string, len(arr))
	for i, v := range arr {
		arr2[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(arr2, "|")
}

func getStringArray(value url.Values, key string) []string {
	str := value.Get(key)

	if str == "" {
		return nil
	}

	return splitByWithEscape(str, '|')
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

func writeTLSRequiredError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("TLS required"))
}

func getPage(values url.Values, key string) int32 {
	v := getInt32FromQuery(values, key)

	if !v.Valid {
		return 1
	}

	if v.Int32 < 1 {
		return 1
	}

	return v.Int32
}
