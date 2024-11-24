package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/joomcode/errorx"
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

func writeRequestError(err error, w http.ResponseWriter) {
	var werr error

	switch err.(type) {
	case jsonBodyError:
		w.WriteHeader(http.StatusBadRequest)
		_, werr = w.Write([]byte(fmt.Sprintf("json body syntax error: %s", err.Error())))
		break
	case *httpError:
		{
			httpErr := err.(*httpError)
			w.WriteHeader(httpErr.StatusCode)
			_, werr = w.Write([]byte(httpErr.Message))
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, werr = w.Write([]byte(fmt.Sprintf("unknown request error: %s", err.Error())))

		break
	}

	if werr != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

func writeUnauthorizedError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte("unauthorized"))
	if err != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

func writeUnprocessableEntity(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, err := w.Write([]byte(message))
	if err != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

func writeApplicationError(w http.ResponseWriter, err error) {
	var werr error

	if errx, ok := err.(*errorx.Error); ok {

		if errorx.HasTrait(errx, app.ErrTraitForbidden) {
			w.WriteHeader(http.StatusForbidden)
			_, werr = w.Write([]byte(err.Error()))
		} else if errorx.HasTrait(errx, app.ErrTraitAuthorizationIssue) {
			w.WriteHeader(http.StatusUnauthorized)
			_, werr = w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusConflict)
			_, werr = w.Write([]byte(err.Error()))
		}

	} else {
		w.WriteHeader(http.StatusConflict)
		_, werr = w.Write([]byte(err.Error()))
	}

	if werr != nil {
		slog.Error("error while writing to the client", "err", err)
	}

}

func write404(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(message))
	if err != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

func writeOK(w http.ResponseWriter) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
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

func int64StringArr(arr []app.Int64String) []int64 {
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

func getStringArray(value url.Values, key string) []string {
	str := value.Get(key)

	if str == "" {
		return nil
	}

	return splitByWithEscape(str, '|')
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
