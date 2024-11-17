package main

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
)

type userController struct {
	service app.UserService
}

func newUserController(service app.UserService) userController {
	return userController{
		service: service,
	}
}

func (c *userController) GetUser(w http.ResponseWriter, r *http.Request) {
	currentUserID := getNullableUserID(r)
	userID, err := urlParamUUID(r, "userID")
	if err != nil {
		writeRequestError(err, w)
		return
	}

	user, err := c.service.GetUser(r.Context(), app.GetUserQuery{
		UserID: currentUserID,
		ID:     userID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, user)
}
