package admin

import (
	"net/http"
	"net/url"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	i18nProvider "github.com/MaratBR/openlibrary/internal/i18n-provider"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/MaratBR/openlibrary/web/olresponse"
	"github.com/ggicci/httpin"
	"github.com/gofrs/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func (c *usersController) Users(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := olhttp.GetPage(query, "p")
	pageSize := olhttp.GetPageSize(query, "ps", 1, 100, 20)
	roles := getUserRoles(query, "usersFilter.role")
	searchQuery := query.Get("q")

	appQuery := app.UsersQuery{
		Page:     page,
		PageSize: pageSize,
		Role:     roles,
		Query:    searchQuery,
	}

	users, err := c.service.ListUsers(r.Context(), appQuery)
	if err != nil {
		olresponse.Write500(w, r, err)
		return
	}

	templates.Users(users, &appQuery).Render(r.Context(), w)
}

func (c *usersController) User(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	c.sendUserEditForm(w, r, userID)
}

func (c *usersController) sendUserEditForm(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	currentUser := auth.RequireUser(r.Context())

	user, err := c.service.GetUserDetails(r.Context(), app.GetUserQuery{
		ID:     userID,
		UserID: uuid.NullUUID{Valid: true, UUID: currentUser.ID},
	})
	if err != nil {
		olresponse.Write500(w, r, err)
		return
	}

	templates.User(user).Render(r.Context(), w)
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
		ActorUserID: uuid.NullUUID{Valid: true, UUID: currentUser.ID},
		UserID:      userID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	l := i18nProvider.GetLocalizer(r.Context())

	flash.Add(r, flash.Text(l.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "admin.users.userWasUpdated",
	})))

	c.sendUserEditForm(w, r, userID)
}
