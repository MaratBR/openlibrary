package server

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
)

type favoritesController struct {
	favoritesService *app.FavoriteService
}

func newFavoritesController(favoritesService *app.FavoriteService) *favoritesController {
	return &favoritesController{favoritesService: favoritesService}
}

func (c *favoritesController) SetFavorite(w http.ResponseWriter, r *http.Request) {
	sessionInfo, ok := getSession(r)
	if !ok {
		// todo implement anonymous favorites
		writeUnauthorizedError(w)
		return
	}

	bookID, err := urlQueryParamInt64(r, "bookId")
	if err != nil {
		writeRequestError(err, w)
		return
	}

	isFavorite := r.URL.Query().Get("isFavorite") == "true"

	err = c.favoritesService.SetFavorite(r.Context(), app.SetFavoriteCommand{
		IsFavorite: isFavorite,
		UserID:     sessionInfo.UserID,
		BookID:     bookID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeOK(w)
}
