package admin

import "net/http"

func writeBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}

func writeApplicationError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}
