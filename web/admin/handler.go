package admin

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/MaratBR/openlibrary/web/webinfra"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var FXModule = fx.Module("http_admin", fx.Provide(
	newTagsController,
	newDebugController,
	newLoginController,
	newUsersController,
	webinfra.AsMountableHandler(newHandler),
))

type Handler struct {
	r chi.Router
}

func newHandler(
	sessionService app.SessionService,
	userService app.UserService,
	loginController *loginController,
	tagsController *tagsController,
	usersController *usersController,
	debugController *debugController,

	flashMiddleware flash.Middleware,
	log *zap.SugaredLogger,
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
		}, log))

		// anonymous area
		r.Get("/login", loginController.Login)

		// authorization required
		r.Group(func(r chi.Router) {
			r.Use(requireAdmin)
			setupSPA(r, usersController, tagsController, debugController)

			// Preserve old form submissions while all GET screens use the SPA.
			r.With(httpin.NewInput(updateUserRequest{})).Post("/users/{id}", usersController.UserUpdate)
			r.With(httpin.NewInput(&tagEditBody{})).Post("/tags/tag-details/{id}/edit", tagsController.TagEdit)

		})
	})
	return h
}

func (h *Handler) MountAt() string { return "/admin" }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func adminNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.NotFound().Render(r.Context(), w)
}
