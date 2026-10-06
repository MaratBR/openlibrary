package admin

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"

	"github.com/ggicci/httpin"
)

type usersController struct {
	service app.UserService
}

func newUsersController(service app.UserService) *usersController {
	return &usersController{service: service}
}

func getUserRoles(query url.Values, key string) []app.UserRole {
	var roles []app.UserRole
	for _, roleStr := range olhttp.GetStringArray(query, key) {
		role, err := app.ParseUserRole(roleStr)
		if err == nil {
			roles = append(roles, role)
		}
	}
	return roles
}

type updateUserRequest struct {
	Gender     string `in:"form=genderOther"`
	GenderType string `in:"form=gender"`
	Password   string `in:"form=password"`
	About      string `in:"form=about"`
	Role       string `in:"form=role"`
}

func (req *updateUserRequest) GetGender() string {
	if req.GenderType == "male" || req.GenderType == "female" || req.GenderType == "" {
		return req.GenderType
	}
	return req.Gender
}
func (req *updateUserRequest) GetRole() app.Nullable[app.UserRole] {
	role, err := app.ParseUserRole(req.Role)
	if err != nil {
		return app.Null[app.UserRole]()
	}
	return app.Value(role)
}

func (c *usersController) UserUpdate(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	currentUser := auth.RequireUser(r.Context())
	input := r.Context().Value(httpin.Input).(*updateUserRequest)

	err = c.service.UpdateUser(r.Context(), app.UpdateUserCommand{
		Password:    input.Password,
		About:       input.About,
		Gender:      input.GetGender(),
		Role:        input.GetRole(),
		ActorUserID: app.Value(currentUser.ID),
		UserID:      userID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	l := i18n.GetLocalizer(r.Context())

	flash.Add(r, flash.Text(l.T("admin.users.userWasUpdated")))

	http.Redirect(w, r, fmt.Sprintf("/admin#/users/%s", userID), http.StatusSeeOther)
}
