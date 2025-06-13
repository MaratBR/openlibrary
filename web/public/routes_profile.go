package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
)

type userController struct {
	service     app.UserService
	bookService app.BookService
}

func newProfileController(service app.UserService, bookService app.BookService) *userController {
	return &userController{service: service, bookService: bookService}
}

func (c *userController) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	user, err := c.service.GetUserDetails(r.Context(), app.GetUserQuery{
		ID: userID,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	pinnedBooks, err := c.bookService.GetPinnedBooks(r.Context(), app.GetPinnedUserBooksQuery{
		UserID: user.ID,
		Limit:  6,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.User(user, pinnedBooks).Render(r.Context(), w)
}
