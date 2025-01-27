package admin

import (
	"net/http"
	"net/url"

	"github.com/MaratBR/openlibrary/internal/app"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/MaratBR/openlibrary/web/olresponse"
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

	users, err := c.service.ListUsers(r.Context(), app.UsersQuery{
		Page:     page,
		PageSize: pageSize,
		Role:     roles,
	})
	if err != nil {
		olresponse.Write500(w, r, err)
		return
	}

	templates.Users(users).Render(r.Context(), w)
}
