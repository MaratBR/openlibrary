package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
)

func adminBookOverride(r *http.Request) bool {
	user, ok := auth.GetUser(r.Context())
	return ok && user.Role.IsAdmin() && r.URL.Query().Get("admin.override") == "1"
}

// Offer an explicit opt-in without reading or displaying restricted content.
func offerAdminBookOverride(w http.ResponseWriter, r *http.Request) bool {
	user, ok := auth.GetUser(r.Context())
	if !ok || !user.Role.IsAdmin() || adminBookOverride(r) {
		return false
	}
	target := *r.URL
	query := target.Query()
	query.Set("admin.override", "1")
	target.RawQuery = query.Encode()
	olhttp.WriteTemplate(w, r.Context(), templates.AdminBookOverridePrompt(target.RequestURI()))
	return true
}
