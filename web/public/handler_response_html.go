package public

import (
	"net/http"
)

func write500(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func writeBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(400)
	if err == nil {
		w.Write([]byte("Something went wrong")) // aka "I am lazy"
	} else {
		w.Write([]byte(err.Error()))

	}
}

func writeApplicationError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(409)
	w.Write([]byte(err.Error()))
}

func writeUnauthorizedError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}
