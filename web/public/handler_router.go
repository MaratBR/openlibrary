package public

import (
	"errors"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/upload"
	"github.com/MaratBR/openlibrary/web/olresponse"
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

	bookManagerService := app.NewBookManagerService(db, tagsService, h.uploadService, userService, bgServices.BookReindex)

	h.r.Group(func(r chi.Router) {
		r.Use(auth.NewAuthorizationMiddleware(sessionService, userService, auth.MiddlewareOptions{
			OnFail: func(w http.ResponseWriter, r *http.Request, err error) {
				olresponse.Write500(w, r, err)
			},
		}))

		r.NotFound(notFoundHandler)
		r.MethodNotAllowed(methodNotAllowed)

		newAuthController(authService, h.csrfHandler).Register(r)
		newBookController(bookService, reviewsService, readingListService).Register(r)
		newChaptersController(bookService).Register(r)
		newSearchController(searchService, bookService).Register(r)
		newTagsController(tagsService).Register(r)
		newProfileController(userService, bookService).Register(r)
		newLibraryController(readingListService).Register(r)
		newBookManagerController(bookManagerService).Register(r)

		r.Route("/mod", func(r chi.Router) {
			newModController(bookService, modBookService).Register(r)
		})

		r.Route("/debug", func(r chi.Router) {
			r.Handle("/500", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				olresponse.Write500(w, r, errors.New("test error"))
			}))
		})

		r.Route("/_api", func(r chi.Router) {
			newAPIBookController(bookService, reviewsService, readingListService).Register(r)
			newAPIReadingListController(readingListService).Register(r)
			newAPITagsController(tagsService).Register(r)
			newAPIBookManagerController(bookManagerService, fileValidator).Register(r)
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
