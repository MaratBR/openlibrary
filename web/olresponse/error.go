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
		writeErrorJSON(err, w)
	} else {
		err500(err).Render(r.Context(), w)
	}
}
