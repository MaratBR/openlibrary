package publicui

import (
	"context"
	_ "embed"
	"net/http"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/csrf"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public-ui/templates"
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
)

type Handler struct {
	_mutex   sync.Mutex
	_started bool

	backgroundServices []interface {
		Start() error
		Stop()
	}
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
	h.initRouter()
	h.setupRouter()
	return h
}

func (h *Handler) initRouter() {
	h.r = chi.NewRouter()
	h.r.Use(gziphandler.GzipHandler)
	h.r.Use(olhttp.ReqCtxMiddleware)

	// add version of the app as context info
	h.r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "version", h.version))
			next.ServeHTTP(w, r)
		})
	})
	h.r.NotFound(notFoundHandler)
}

// Start starts all background services and sets started flag to true.
// If any of the services fail to start, an error is returned.
func (h *Handler) Start() error {
	h._mutex.Lock()
	defer h._mutex.Unlock()
	if h._started {
		return nil
	}

	for _, s := range h.backgroundServices {
		err := s.Start()
		if err != nil {
			return err
		}
	}

	h._started = true
	return nil
}

func (h *Handler) Stop() {
	h._mutex.Lock()
	defer h._mutex.Unlock()
	if !h._started {
		return
	}

	for _, s := range h.backgroundServices {
		s.Stop()
	}

	h._started = false
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h._started {
		panic("cannot serve http until handle has been started")
	}

	h.r.ServeHTTP(w, r)
}

var (
	Dev = true
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	templates.NotFoundPage().Render(r.Context(), w)
}
