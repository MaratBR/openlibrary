package public

import (
	"errors"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/upload"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) setupRouter(bgServices *app.BackgroundServices) {
	db := h.db

	fileValidator := upload.NewFileValidator(h.cfg)

	// application layer services
	sessionService := app.NewCachedSessionService(app.NewSessionService(db), h.cache)
	authService := app.NewAuthService(db, sessionService)

	tagsService := app.NewTagsService(db)
	readingListService := app.NewReadingListService(db, h.uploadService)
	userService := app.NewUserService(db)
	reviewsService := app.NewCachedReviewsService(app.NewReviewsService(db, userService, bgServices.Book), h.cache)
	bookService := app.NewBookService(db, tagsService, h.uploadService, readingListService, reviewsService)
	modBookService := app.NewModerationBookService(db)
	searchService := app.NewCachedSearchService(app.NewSearchService(db, tagsService, h.uploadService, userService, h.esClient), h.cache)
	collectionService := app.NewCollectionsService(db, tagsService, h.uploadService)

	bookManagerService := app.NewBookManagerService(db, tagsService, h.uploadService, userService, bgServices.BookReindex)

	h.r.Group(func(r chi.Router) {
		r.Use(auth.NewAuthorizationMiddleware(sessionService, userService, auth.MiddlewareOptions{
			OnFail: func(w http.ResponseWriter, r *http.Request, err error) {
				olhttp.Write500(w, r, err)
			},
		}))

		r.NotFound(notFoundHandler)
		r.MethodNotAllowed(methodNotAllowed)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			olhttp.WriteTemplate(w, r.Context(), templates.Home())
		})

		newAuthController(authService, h.csrfHandler).Register(r)
		newBookController(bookService, reviewsService, readingListService).Register(r)
		newChaptersController(bookService, readingListService).Register(r)
		newSearchController(searchService, bookService).Register(r)
		newTagsController(tagsService).Register(r)
		newProfileController(userService, bookService, searchService, collectionService).Register(r)
		newLibraryController(readingListService, collectionService).Register(r)
		newCollectionController(collectionService).Register(r)
		newBookManagerController(bookManagerService, collectionService).Register(r)

		r.Route("/mod", func(r chi.Router) {
			newModController(bookService, modBookService).Register(r)
		})

		r.Route("/debug", func(r chi.Router) {
			r.Handle("/500", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				olhttp.Write500(w, r, errors.New("test error"))
			}))
		})

		r.Route("/_api", func(r chi.Router) {
			newAPIBookController(bookService, reviewsService, readingListService).Register(r)
			newAPIReadingListController(readingListService).Register(r)
			newAPITagsController(tagsService).Register(r)
			newAPIBookManagerController(bookManagerService, fileValidator).Register(r)
			newAPICollectionController(collectionService).Register(r)

			r.NotFound(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				olhttp.NewAPIError(errors.New("not found")).Write(w)
			})
		})

	})
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func redirectToLoginOnUnauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.GetSession(r.Context())

		if !ok {
			redirectToLogin(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
