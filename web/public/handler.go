package public

import (
	_ "embed"
	"net/http"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/NYTimes/gziphandler"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	_mutex   sync.Mutex
	_started bool

	r             chi.Router
	db            app.DB
	esClient      *elasticsearch.TypedClient
	cfg           *koanf.Koanf
	redisClient   *redis.Client
	cache         *cache.Cache
	csrfHandler   *csrf.Handler
	uploadService *app.UploadService
}

func NewHandler(
	db app.DB,
	cfg *koanf.Koanf,
	redisService *redis.Client,
	cache *cache.Cache,
	csrfHandler *csrf.Handler,
	bgServices *app.BackgroundServices,
	uploadService *app.UploadService,
	esClient *elasticsearch.TypedClient,
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
	if esClient == nil {
		panic("esClient is nil")
	}

	h := &Handler{
		db:            db,
		cfg:           cfg,
		cache:         cache,
		redisClient:   redisService,
		csrfHandler:   csrfHandler,
		uploadService: uploadService,
		esClient:      esClient,
	}
	h.initRouter()
	h.setupRouter(bgServices)
	return h
}

func (h *Handler) initRouter() {
	h.r = chi.NewRouter()
	h.r.Use(gziphandler.GzipHandler)
	h.r.Use(flash.Middleware)

	h.r.NotFound(notFoundHandler)
	h.r.MethodNotAllowed(methodNotAllowed)
}

// Start starts all background services and sets started flag to true.
// If any of the services fail to start, an error is returned.
func (h *Handler) Start() error {
	h._mutex.Lock()
	defer h._mutex.Unlock()
	if h._started {
		return nil
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

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)

	if olhttp.PreferredMimeTypeIsJSON(r) {
		w.Write([]byte("Method Not Allowed"))
	} else {
		templates.MethodNotAllowedPage().Render(r.Context(), w)
	}
}
