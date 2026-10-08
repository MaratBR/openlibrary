package public

import (
	"net/http"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/store"
)

func bypassContentPreferences(r *http.Request) bool {
	for _, prefix := range []string{"/admin", "/books-manager", "/_api/books-manager", "/moderation", "/_api/moderation", "/mod", "/account/settings"} {
		if r.URL.Path == prefix || strings.HasPrefix(r.URL.Path, prefix+"/") {
			return true
		}
	}
	user, ok := auth.GetUser(r.Context())
	return ok && user.Role.IsAdmin() && (r.URL.Query().Get("admin.link") == "1" || r.URL.Query().Get("admin.override") == "1")
}

func contentPreferencesMiddleware(users app.UserService, db store.DBTX) func(http.Handler) http.Handler {
	queries := store.New(db)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := auth.GetUser(r.Context())
			if !ok || bypassContentPreferences(r) {
				next.ServeHTTP(w, r)
				return
			}
			settings, err := users.GetUserModerationSettings(r.Context(), user.ID)
			if err != nil {
				olhttp.Write500(w, r, err)
				return
			}
			prefs, err := app.ResolveContentPreferences(r.Context(), queries, settings)
			if err != nil {
				olhttp.Write500(w, r, err)
				return
			}
			// Personalized responses and blur decisions must never enter shared caches.
			w.Header().Set("Cache-Control", "private, no-store")
			next.ServeHTTP(w, r.WithContext(app.WithContentPreferences(r.Context(), prefs)))
		})
	}
}
