package public

import (
	"log/slog"
	"net/http"

	"github.com/joomcode/errorx"
)

func apiWrite500(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func apiWriteBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}

func apiWriteApplicationError(w http.ResponseWriter, err error) {
	w.WriteHeader(409)
	w.Write([]byte(err.Error()))
}

func apiWriteUnexpectedApplicationError(w http.ResponseWriter, err error) {
	var werr error
	w.WriteHeader(http.StatusInternalServerError)

	if errx, ok := err.(*errorx.Error); ok {
		_, werr = w.Write([]byte(errx.Error()))
	} else {
		_, werr = w.Write([]byte(err.Error()))
	}

	if werr != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}

func apiWriteUnprocessableEntity(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, err := w.Write([]byte(message))
	if err != nil {
		slog.Error("error while writing to the client", "err", err)
	}
}
