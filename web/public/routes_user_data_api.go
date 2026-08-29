package public

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/go-chi/chi/v5"
)

type apiControllerUserData struct {
	service app.UserDataService
}

func newAPIUserDataController(service app.UserDataService) *apiControllerUserData {
	return &apiControllerUserData{service: service}
}

func (c *apiControllerUserData) Register(r chi.Router) {
	r.With(apiRequiresAuthorizationMiddleware).Get("/user-data/{key}", c.get)
	r.With(apiRequiresAuthorizationMiddleware).Put("/user-data/{key}", c.set)
}

func (c *apiControllerUserData) get(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	data, err := c.service.Get(r.Context(), app.GetUserDataQuery{
		UserID: session.UserID,
		Key:    chi.URLParam(r, "key"),
	})
	if err != nil {
		c.writeError(w, err)
		return
	}
	olhttp.NewAPIResponse(data).Write(w)
}

func (c *apiControllerUserData) set(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	var data json.RawMessage
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, app.UserDataMaxSize+1))
	if err := decoder.Decode(&data); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			err = errors.New("request body must contain a single JSON value")
		}
		apiWriteBadRequest(w, err)
		return
	}

	err := c.service.Set(r.Context(), app.SetUserDataQuery{
		UserID: session.UserID,
		Key:    chi.URLParam(r, "key"),
		Data:   data,
	})
	if err != nil {
		c.writeError(w, err)
		return
	}
	apiWriteOK(w)
}

func (c *apiControllerUserData) writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, app.ErrUserDataKeyNotAllowed) || errors.Is(err, app.ErrUserDataTooLarge) {
		apiWriteUnprocessableEntity(w, err)
		return
	}
	apiWriteUnexpectedApplicationError(w, err)
}
