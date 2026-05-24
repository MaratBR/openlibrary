package public

import (
	"net/http"
)

const canViewAdultContentCookieName string = "view_adult"

func canViewAdultContent(r *http.Request) bool {
	cookie, err := r.Cookie(canViewAdultContentCookieName)
	if err == nil && cookie.Value == "1" {
		return true
	}

	// TODO return this eventually
	// user, ok := auth.GetUser(r.Context())
	// if ok && user.ShowAdultContent {
	// 	return true
	// }

	return false
}
