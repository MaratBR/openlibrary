package admin

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) setupRouter(bgServices *app.BackgroundServices) {
	h.r.NotFound(adminNotFound)

	h.r.Group(func(r chi.Router) {
		sessionService := app.NewCachedSessionService(app.NewSessionService(h.db), h.cache)
		tagsService := app.NewTagsService(h.db)
		userService := app.NewUserService(h.db)

		r.Use(flash.Middleware)
		r.Use(auth.NewAuthorizationMiddleware(sessionService, userService, auth.MiddlewareOptions{
			OnFail: func(w http.ResponseWriter, r *http.Request, err error) {
				olhttp.Write500(w, r, err)
			},
		}))

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
				newTagsController(h.db, h.cfg, tagsService).Setup(r)
			})

			r.Route("/users", func(r chi.Router) {
				c := newUsersController(userService)

				r.Get("/", c.Users)
				r.Get("/{id}", c.User)
				r.With(httpin.NewInput(updateUserRequest{})).Post("/{id}", c.UserUpdate)
			})

			r.Route("/books", func(r chi.Router) {
				newBooksController(h.db).Register(r)
			})

			{
				c := newDebugController(bgServices.BookReindex)
				r.Handle("/debug", http.HandlerFunc(c.Actions))
			}
		})
	})

}
