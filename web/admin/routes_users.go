package admin

import (
	"net/http"

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

func (c *usersController) Users(w http.ResponseWriter, r *http.Request) {

	page := olhttp.GetPage(r.URL.Query(), "p")
	pageSize := olhttp.GetPageSize(r.URL.Query(), "ps", 1, 100, 20)

	users, err := c.service.ListUsers(r.Context(), app.UsersQuery{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		olresponse.Write500(w, r, err)
		return
	}

	templates.Users(users).Render(r.Context(), w)
}
