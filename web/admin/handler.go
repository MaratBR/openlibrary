package admin

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/MaratBR/openlibrary/web/webfx"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"
)

var FXModule = fx.Module("http_admin", fx.Provide(
	newTagsController,
	newBooksController,
	newDebugController,
	newLoginController,
	newUsersController,
	webfx.AsMountableHandler(newHandler),
))

type Handler struct {
	r chi.Router
}

func newHandler(
	db app.DB,
	cfg *koanf.Koanf,
	cache *cache.Cache,
	sessionService app.SessionService,
	userService app.UserService,
	esClient *elasticsearch.TypedClient,
	loginController *loginController,
	tagsController *tagsController,
	usersController *usersController,
	debugController *debugController,
	booksController *booksController,

	flashMiddleware flash.Middleware,
) *Handler {
	h := &Handler{
		r: chi.NewRouter(),
	}
	h.r.NotFound(adminNotFound)

	h.r.Group(func(r chi.Router) {
		r.Use(flashMiddleware)
		r.Use(auth.NewAuthorizationMiddleware(sessionService, userService, auth.MiddlewareOptions{
			OnFail: func(w http.ResponseWriter, r *http.Request, err error) {
				olhttp.Write500(w, r, err)
			},
		}))

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
				tagsController.Setup(r)
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/", usersController.Users)
				r.Get("/{id}", usersController.User)
				r.With(httpin.NewInput(updateUserRequest{})).Post("/{id}", usersController.UserUpdate)
			})

			r.Route("/books", func(r chi.Router) {
				booksController.Register(r)
			})

			r.Handle("/debug", http.HandlerFunc(debugController.Actions))
		})
	})
	return h
}

func (h *Handler) MountAt() string { return "/admin" }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func adminNotFound(w http.ResponseWriter, r *http.Request) {
	templates.NotFound().Render(r.Context(), w)
}
