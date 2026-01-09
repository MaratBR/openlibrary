package public

import (
	_ "embed"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/MaratBR/openlibrary/web/webfx"
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type Handler struct {
	r chi.Router
}

var FXModule = fx.Module("public_ui_handler",

	fx.Provide(
		newHomeController,
		newAuthController,
		newBookController,
		newChaptersController,
		newSearchController,
		newTagsController,
		newProfileController,
		newLibraryController,
		newCollectionController,
		newBookManagerController,
		newModController,
		newAPIBookController,
		newAPIReadingListController,
		newAPITagsController,
		newAPIBookManagerController,
		newAPICollectionController,
		newAPICommentsController,
		webfx.AsMountableHandler(newHandler),
	))

// i should really do something with number of params in this function...

func newHandler(
	sessionService app.SessionService,
	userService app.UserService,

	homeController *homeController,
	authController *authController,
	bookController *bookController,
	bookManagerController *bookManagerController,
	chapterController *chaptersController,
	collectionController *collectionController,
	modController *modController,
	libraryController *libraryController,
	profileController *profileController,
	searchController *searchController,
	tagsController *tagsController,

	apiControllerTags *apiControllerTags,
	apiControllerBook *apiControllerBook,
	apiControllerBookManager *apiControllerBookManager,
	apiControllerCollection *apiControllerCollection,
	apiControllerReadingList *apiControllerReadingList,
	apiControllerComments *apiControllerComments,

	flashMiddleware flash.Middleware,
) webfx.MountableHandler {
	h := &Handler{}

	h.r = chi.NewRouter()
	h.r.Use(gziphandler.GzipHandler)
	h.r.Use(flashMiddleware)

	h.r.NotFound(notFoundHandler)
	h.r.MethodNotAllowed(methodNotAllowed)

	h.r.Use(auth.NewAuthorizationMiddleware(sessionService, userService, auth.MiddlewareOptions{
		OnFail: func(w http.ResponseWriter, r *http.Request, err error) {
			olhttp.Write500(w, r, err)
		},
	}))

	homeController.Register(h.r)
	authController.Register(h.r)
	bookController.Register(h.r)
	bookManagerController.Register(h.r)
	chapterController.Register(h.r)
	collectionController.Register(h.r)
	modController.Register(h.r)
	libraryController.Register(h.r)
	profileController.Register(h.r)
	searchController.Register(h.r)
	tagsController.Register(h.r)

	h.r.Route("/debug", func(r chi.Router) {
		r.Handle("/500", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			olhttp.Write500(w, r, errors.New("test error"))
		}))
	})

	h.r.Route("/_api", func(r chi.Router) {

		apiControllerBook.Register(r)
		apiControllerBookManager.Register(r)
		apiControllerCollection.Register(r)
		apiControllerReadingList.Register(r)
		apiControllerTags.Register(r)
		apiControllerComments.Register(r)

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			olhttp.NewAPIError(errors.New("not found")).Write(w)
		})
	})

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func (h *Handler) MountAt() string {
	return "/"
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	templates.NotFoundPage().Render(r.Context(), w)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)

	if olhttp.PreferredMimeTypeIsJSON(r) {
		w.Write([]byte("Method Not Allowed"))
	} else {
		templates.MethodNotAllowedPage().Render(r.Context(), w)
	}
}

func redirectWithNextParameter(w http.ResponseWriter, r *http.Request, path string) {
	next := r.URL.Path
	if r.URL.RawQuery != "" {
		next += "?" + r.URL.RawQuery
	}

	u, err := url.Parse(path)
	if err != nil {
		slog.Error("failed to parse redirect url")
	} else {
		// TODO remove next param is next is the same as the URL we are redirecting too, for some reason
		q := u.Query()
		q.Set("next", next)
		u.RawQuery = q.Encode()
		path = u.String()
	}

	http.Redirect(w, r, path, http.StatusFound)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	redirectWithNextParameter(w, r, "/login")
}

func requiresAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.GetUser(r.Context())

		if !ok {
			redirectToLogin(w, r)
			return
		}

		if !user.IsEmailVerified {
			redirectWithNextParameter(w, r, "/signup/email-verification-code")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func apiRequiresAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.GetSession(r.Context())

		if !ok {
			apiWriteUnauthorized(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
