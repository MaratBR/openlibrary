package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/auth"
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

func apiWriteUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	resp := olresponse.NewAPIErrorWithMessage("unauthorized")
	resp.Write(w)
}

// I added ____ because for some reason autocomplete REALLY loves to suggest this function whenever I want to type "return"
// TODO rename
func ___returnUnauthorizedIfNotLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.GetSession(r.Context())
		if !ok {
			apiWriteUnauthorized(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
