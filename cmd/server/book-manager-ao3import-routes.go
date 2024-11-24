package server

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
)

type ao3ImportRequest struct {
	ID string `json:"id"`
}

type ao3ImportResponse struct {
	ID int64 `json:"id,string"`
}

func (b *bookManagerController) ImportAO3(w http.ResponseWriter, r *http.Request) {
	sessionInfo := requireSession(r)
	req, err := getJSON[ao3ImportRequest](r)
	if err != nil {
		writeRequestError(err, w)
		return
	}
	bookID, err := b.service.ImportFromBookAo3(r.Context(), app.ManagerCreateBookFromAo3Command{
		UserID: sessionInfo.UserID,
		Ao3ID:  req.ID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, ao3ImportResponse{
		ID: bookID,
	})
}
