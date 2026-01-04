package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/auth"
)

const canViewAdultContentCookieName string = "view_adult"

func canViewAdultContent(r *http.Request) bool {
	cookie, err := r.Cookie(canViewAdultContentCookieName)
	if err == nil && cookie.Value == "1" {
		return true
	}

	user, ok := auth.GetUser(r.Context())
	if ok && user.ShowAdultContent {
		return true
	}

	return false
}
