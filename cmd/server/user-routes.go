package server

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

	user, err := c.service.GetUserDetails(r.Context(), app.GetUserQuery{
		UserID: currentUserID,
		ID:     userID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, user)
}

type whoamiResponse struct {
	User *app.SelfUserDto `json:"user"`
}

func (c *userController) Whoami(w http.ResponseWriter, r *http.Request) {
	var response whoamiResponse

	currentUserID := getNullableUserID(r)
	if !currentUserID.Valid {
		writeJSON(w, response)
		return
	}

	user, err := c.service.GetUserSelfData(r.Context(), currentUserID.UUID)
	if err != nil {
		writeApplicationError(w, err)
		return
	} else {
		response.User = user
	}
	writeJSON(w, response)
}

func (c *userController) Follow(w http.ResponseWriter, r *http.Request) {
	userID, err := urlQueryParamUUID(r, "userId")
	if err != nil {
		writeRequestError(err, w)
		return
	}
	session := requireSession(r)
	err = c.service.FollowUser(r.Context(), app.FollowUserCommand{
		UserID:   userID,
		Follower: session.UserID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeOK(w)
}

func (c *userController) Unfollow(w http.ResponseWriter, r *http.Request) {
	userID, err := urlQueryParamUUID(r, "userId")
	if err != nil {
		writeRequestError(err, w)
		return
	}
	session := requireSession(r)
	err = c.service.UnfollowUser(r.Context(), app.UnfollowUserCommand{
		UserID:   userID,
		Follower: session.UserID,
	})
	if err != nil {
		writeApplicationError(w, err)
	}
	writeOK(w)
}
