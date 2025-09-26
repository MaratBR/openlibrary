package olresponse

import (
	"net/http"
)

func Write500(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	w.WriteHeader(http.StatusInternalServerError)

	if PreferredMimeTypeIsJSON(r) {
		NewAPIError(err).Write(w)
	} else {
		err500(err).Render(r.Context(), w)
	}
}

func WriteCustomErrorPage(
	w http.ResponseWriter,
	r *http.Request,
	title, subtitle string,
	err error,
) {
	if PreferredMimeTypeIsJSON(r) {
		NewAPIError(err).Write(w)
	} else {
		errorPageWithACat(title, subtitle, err.Error(), -1).Render(r.Context(), w)
	}
}
