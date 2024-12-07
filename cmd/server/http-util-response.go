package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/joomcode/errorx"
)

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
