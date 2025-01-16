package publicui

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/web/public-ui/templates"
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
)

type Handler struct {
	r           chi.Router
	db          app.DB
	cfg         *koanf.Koanf
	cache       *cache.Cache
	csrfHandler *csrf.Handler
	version     string
}

func NewHandler(
	db app.DB,
	cfg *koanf.Koanf,
	version string,
	cache *cache.Cache,
	csrfHandler *csrf.Handler,
) *Handler {
	if cache == nil {
		panic("cache is nil")
	}
	if cfg == nil {
		panic("cfg is nil")
	}
	if db == nil {
		panic("db is nil")
	}

	h := &Handler{
		db:          db,
		cfg:         cfg,
		version:     version,
		cache:       cache,
		csrfHandler: csrfHandler,
	}
	h.createRouter()
	return h
}

func (h *Handler) createRouter() {
	h.r = chi.NewRouter()
	h.r.Use(gziphandler.GzipHandler)
	h.r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "version", h.version))
			next.ServeHTTP(w, r)
		})
	})

	h.r.Mount("/_/embed-assets/", http.StripPrefix("/_/embed-assets/", newAssetsHandler()))

	if Dev {
		h.r.Mount("/_/assets/", http.StripPrefix("/_/assets/", http.FileServer(http.Dir("web/public-ui/frontend/dist"))))
	} else {
		panic("not implemented yet")
	}

	h.r.NotFound(notFoundHandler)

	db := h.db

	// application layer services
	uploadService := app.NewUploadServiceFromApplicationConfig(h.cfg)
	sessionService := app.NewSessionService(db)
	sessionService = app.NewCachedSessionService(sessionService, h.cache)
	authService := app.NewAuthService(db, sessionService)
	// favoriteRecalculationBackgroundService := app.NewFavoriteRecalculationBackgroundService(db)
	// favoritesService := app.NewFavoriteService(db, favoriteRecalculationBackgroundService)
	tagsService := app.NewTagsService(db)
	readingListService := app.NewReadingListService(db)
	userService := app.NewUserService(db)
	// bookManagerService := app.NewBookManagerService(db, tagsService, uploadService)
	bookBackgroundService := app.NewBookBackgroundService(db)
	reviewsService := app.NewReviewsService(db, userService, bookBackgroundService)
	reviewsService = app.NewCachedReviewsService(reviewsService, h.cache)
	bookService := app.NewBookService(db, tagsService, uploadService, readingListService, reviewsService)
	searchService := app.NewSearchService(db, tagsService, uploadService)
	searchService = app.NewCachedSearchService(searchService, h.cache)

	h.r.Group(func(r chi.Router) {
		authController := newAuthController(authService, h.csrfHandler)
		bookController := newBookController(bookService, reviewsService, readingListService)
		chapterController := newChaptersController(bookService)

		r.HandleFunc("/login", authController.LogIn)

		r.Get("/book/{bookID}", bookController.GetBook)

		r.Get("/book/{bookID}/chapters/{chapterID}", chapterController.GetChapter)
	})

	h.r.Route("/_api", func(r chi.Router) {
		apiBookController := newAPIBookController(bookService, reviewsService, readingListService)
		apiReadingListController := newAPIReadingListController(readingListService)

		r.Post("/reviews/rating", apiBookController.RateBook)

		r.Post("/reading-list/status", apiReadingListController.UpdateStatus)
	})

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

var (
	Dev = true
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	templates.NotFoundPage(r.Context()).Render(r.Context(), w)
}
