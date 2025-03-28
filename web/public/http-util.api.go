package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/web/olresponse"
)

func apiWrite500(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	resp := olresponse.NewAPIError(err)
	resp.Write(w)
}

func apiWriteBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	resp := olresponse.NewAPIError(err)
	resp.Write(w)
}

func apiWriteOK(w http.ResponseWriter) {
	w.WriteHeader(200)
	olresponse.NewAPIResponseOK().Write(w)
}

func apiWriteApplicationError(w http.ResponseWriter, err error) {
	w.WriteHeader(409)
	resp := olresponse.NewAPIError(err)
	resp.Write(w)
}

func apiWriteUnexpectedApplicationError(w http.ResponseWriter, err error) {
	apiWrite500(w, err)
}

func apiWriteUnprocessableEntity(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	resp := olresponse.NewAPIError(err)
	resp.Write(w)
}
