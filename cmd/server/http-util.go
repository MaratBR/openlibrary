package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/go-chi/chi/v5"
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
	switch err.(type) {
	case jsonBodyError:
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(fmt.Sprintf("json body syntax error: %s", err.Error())))
		if err != nil {
			slog.Error("error while writing to the client", "err", err)
		}
		break
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(fmt.Sprintf("unknown request error: %s", err.Error())))
		if err != nil {
			slog.Error("error while writing to the client", "err", err)
		}
		break
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

func int64StringArr(arr []app.Int64String) []int64 {
	arr2 := make([]int64, len(arr))
	for i, v := range arr {
		arr2[i] = int64(v)
	}
	return arr2
}
