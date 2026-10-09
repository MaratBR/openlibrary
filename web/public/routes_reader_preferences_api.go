package public

import (
	"encoding/json"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/go-chi/chi/v5"
)

type apiControllerReaderPreferences struct {
	service app.ReaderPreferencesService
	fonts   app.ReaderFontService
}

func newAPIReaderPreferencesController(service app.ReaderPreferencesService, fonts app.ReaderFontService) *apiControllerReaderPreferences {
	return &apiControllerReaderPreferences{service: service, fonts: fonts}
}

func (c *apiControllerReaderPreferences) Register(r chi.Router) {
	r.Post("/reader-preferences", c.save)
}

func (c *apiControllerReaderPreferences) save(w http.ResponseWriter, r *http.Request) {
	session, ok := auth.GetSession(r.Context())
	if !ok {
		apiWriteUnauthorized(w)
		return
	}

	var preferences app.ReaderPreferences
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&preferences); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	if err := preferences.ValidateWithFonts(c.fonts.List()); err != nil {
		apiWriteUnprocessableEntity(w, err)
		return
	}
	if err := c.service.Save(r.Context(), session.UserID, preferences); err != nil {
		apiWriteUnexpectedApplicationError(w, err)
		return
	}
	apiWriteOK(w)
}
