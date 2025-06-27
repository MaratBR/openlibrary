package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/reqid"
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
	searchService := app.NewCachedSearchService(app.NewSearchService(db, tagsService, h.uploadService, userService, h.esClient), h.cache)

	bookManagerService := app.NewBookManagerService(db, tagsService, h.uploadService, userService, bgServices.BookReindex)

	h.r.Group(func(r chi.Router) {
		r.Use(auth.NewAuthorizationMiddleware(sessionService, userService, auth.MiddlewareOptions{
			OnFail: func(w http.ResponseWriter, r *http.Request, err error) {
				olresponse.Write500(w, r, err)
			},
		}))
		r.Use(reqid.New())

		authController := newAuthController(authService, h.csrfHandler)
		bookController := newBookController(bookService, reviewsService, readingListService)
		chapterController := newChaptersController(bookService)
		searchController := newSearchController(searchService, bookService)
		tagsController := newTagsController(tagsService)

		// public auth pages
		r.HandleFunc("/login", authController.LogIn)
		r.HandleFunc("/logout", authController.LogOut)

		// book page and its fragments
		r.Get("/book/{bookID}", bookController.GetBook)
		r.Get("/book/{bookID}/__fragment/preview-card", bookController.GetBookPreview)
		r.Get("/book/{bookID}/__fragment/toc", bookController.GetBookTOC)
		r.Get("/book/{bookID}/__fragment/review", bookController.GetBookReview)

		// chapters page
		r.Get("/book/{bookID}/chapters/{chapterID}", chapterController.GetChapter)

		// search controller
		searchController.Register(r)

		r.Get("/tag/{tagID}", tagsController.TagPage)

		r.Route("/users", func(r chi.Router) {
			c := newProfileController(userService, bookService)
			r.Get("/{id}", c.GetProfile)
		})

		r.Route("/library", func(r chi.Router) {
			c := newLibraryController(readingListService)
			c.Register(r)
		})

		r.Route("/_api", func(r chi.Router) {
			apiBookController := newAPIBookController(bookService, reviewsService, readingListService)
			apiReadingListController := newAPIReadingListController(readingListService)
			apiTagsController := newAPITagsController(tagsService)

			r.Post("/reviews/rating", apiBookController.RateBook)
			r.Post("/reviews/{bookID}", apiBookController.UpdateOrCreateReview)
			r.Delete("/reviews/{bookID}", apiBookController.DeleteReview)

			r.Post("/reading-list/status", apiReadingListController.UpdateStatus)

			r.Get("/tags", apiTagsController.Tags)
		})

		newBookManagerController(bookManagerService).Register(r)
		newApiBookManagerController(bookManagerService, fileValidator).Register(r)
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
