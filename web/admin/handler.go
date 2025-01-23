package admin

import (
	"net/http"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/admin/templates"
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
	db    app.DB
	cfg   *koanf.Koanf
	r     chi.Router
	cache *cache.Cache
}

func NewHandler(db app.DB, cfg *koanf.Koanf, cache *cache.Cache) *Handler {
	h := &Handler{db: db, cfg: cfg}
	h.initRouter()
	return h
}

func (h *Handler) initRouter() {
	h.r = chi.NewRouter()
	h.r.Use(olhttp.ReqCtxMiddleware)
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

func adminNotFound(w http.ResponseWriter, r *http.Request) {
	templates.NotFound().Render(r.Context(), w)
}
