package admin

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/url"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/web/admin/templates"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
)

type Handler struct {
	db  app.DB
	cfg *koanf.Koanf
	r   chi.Router
}

func NewHandler(db app.DB, cfg *koanf.Koanf) *Handler {
	h := &Handler{db: db, cfg: cfg}
	h.createRouter()
	return h
}

func (h *Handler) createRouter() {
	h.r = chi.NewRouter()

	adminOrigin := h.cfg.String("server.public-admin-origin")
	if adminOrigin != "" {
		adminOriginU, err := url.Parse(adminOrigin)
		if err != nil {
			panic(err)
		}

		h.r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Host != adminOriginU.Host {
					location := fmt.Sprintf("%s://%s/admin", adminOriginU.Scheme, adminOriginU.Host)
					w.Header().Add("Location", location)
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(307)
					w.Write([]byte(fmt.Sprintf("Redirect to <a href=\"%s\">%s</a>", location, location)))
					return
				}
				next.ServeHTTP(w, r)
			})
		})
	}

	h.r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				session, ok := auth.GetSession(r.Context())
				if ok && session.UserRole == app.RoleAdmin {
					next.ServeHTTP(w, r)
				} else {
					http.Redirect(w, r, "/admin/403", http.StatusFound)
				}
			})
		})

		r.Route("/tags", func(r chi.Router) {
			c := newTagsController(h.db, h.cfg)

			r.Get("/", c.Home)
		})

	})

	h.r.NotFound(adminNotFound)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func adminNotFound(w http.ResponseWriter, r *http.Request) {

	_, ok := auth.GetSession(r.Context())
	if ok {
		templates.NotFound(r.Context()).Render(r.Context(), w)
		return
	} else {
		if r.URL.Path == "/admin" || r.URL.Path == "/admin/login" {
			templates.NotFound(r.Context()).Render(r.Context(), w)
		} else {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
		}
	}
}
