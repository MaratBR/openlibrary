package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type moderationPortalController struct {
	authorizer app.ModerationAuthorizer
}

func newModerationPortalController(authorizer app.ModerationAuthorizer) *moderationPortalController {
	return &moderationPortalController{authorizer: authorizer}
}

func (c *moderationPortalController) Register(r chi.Router) {
	r.Route("/moderation", func(r chi.Router) {
		r.Use(requiresAuthorizationMiddleware)
		r.Use(c.requiresModerator)

		r.Get("/", c.index)
	})
}

func (c *moderationPortalController) requiresModerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := auth.RequireSession(r.Context())
		if err := c.authorizer.AuthorizeModerator(r.Context(), session.UserID); err != nil {
			writeApplicationError(w, r, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *moderationPortalController) index(w http.ResponseWriter, r *http.Request) {
	templates.ModerationPortal().Render(r.Context(), w)
}
