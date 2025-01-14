package publicui

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/public-ui/templates"
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
)

type Handler struct {
	r       chi.Router
	db      app.DB
	cfg     *koanf.Koanf
	cache   *cache.Cache
	version string
}

func NewHandler(db app.DB, cfg *koanf.Koanf, version string) *Handler {
	h := &Handler{
		db:      db,
		cfg:     cfg,
		version: version,
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

	if Dev {
		h.r.Mount("/_/assets/", http.StripPrefix("/_/assets/", http.FileServer(http.Dir("internal/public-ui/frontend/dist"))))
	} else {
		panic("not implemented yet")
	}

	h.r.NotFound(notFoundHandler)

	db := h.db

	// application layer services
	uploadService := app.NewUploadServiceFromApplicationConfig(h.cfg)
	sessionService := app.NewSessionService(db)
	sessionService = app.NewCachedSessionService(sessionService, h.cache)
	// authService := app.NewAuthService(db, sessionService)
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

	{
		bookController := newBookController(bookService, reviewsService)
		h.r.Get("/book/{bookID}", bookController.GetBook)
	}

	{
		chapterController := newChaptersController(bookService)
		h.r.Get("/book/{bookID}/chapters/{chapterID}", chapterController.GetChapter)
	}

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

var (
	Dev = true
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	templates.NotFoundPage(r.Context()).Render(r.Context(), w)
}
