package admin

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/reqid"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/MaratBR/openlibrary/web/olresponse"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) setupRouter(bgServices *app.BackgroundServices) {
	h.setupAutoRedirect()

	h.r.Group(func(r chi.Router) {
		sessionService := app.NewCachedSessionService(app.NewSessionService(h.db), h.cache)
		tagsService := app.NewTagsService(h.db)
		userService := app.NewUserService(h.db)

		authorizationMiddleware := auth.NewAuthorizationMiddleware(sessionService, userService, auth.MiddlewareOptions{
			OnFail: func(w http.ResponseWriter, r *http.Request, err error) {
				olresponse.Write500(w, r, err)
			},
		})

		r.Use(reqid.New())
		r.Use(authorizationMiddleware)

		// controllers for anonymous area
		loginController := newLoginController()

		// anonymous area
		h.r.HandleFunc("/login", loginController.Login)

		// authorization required
		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					session, ok := auth.GetSession(r.Context())
					if !ok {
						http.Redirect(w, r, "/admin/login", http.StatusFound)
						return
					}
					if session.UserRole != app.RoleAdmin {
						w.WriteHeader(http.StatusForbidden)
						templates.Forbidden().Render(r.Context(), w)
						return
					}
					next.ServeHTTP(w, r)
				})
			})

			r.Route("/tags", func(r chi.Router) {
				c := newTagsController(h.db, h.cfg, tagsService)

				r.Get("/", c.Home)
				r.Get("/tag-details/{id}", c.Tag)

			})

			r.NotFound(adminNotFound)
		})
	})

}

func (h *Handler) setupAutoRedirect() {
	adminOrigin := h.cfg.String("server.public-admin-origin")
	if adminOrigin != "" {
		adminOriginU, err := url.Parse(adminOrigin)
		if err != nil {
			panic(err)
		}

		h.r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Host != adminOriginU.Host {
					location := fmt.Sprintf("%s://%s/admin", adminOriginU.Scheme, adminOriginU.Host)
					w.Header().Add("Location", location)
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(307)
					w.Write([]byte(fmt.Sprintf("Redirect to <a href=\"%s\">%s</a>", location, location)))
					return
				}
				next.ServeHTTP(w, r)
			})
		})
	}
}
