package publicui

import "net/http"

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
