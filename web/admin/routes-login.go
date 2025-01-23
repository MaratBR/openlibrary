package admin

import (
	"net/http"
)

type loginController struct{}

func newLoginController() *loginController {
	return &loginController{}
}

func (c *loginController) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login?next=/admin&ol.from=ADMIN_LOGIN", http.StatusFound)
}
