package public

import (
	"errors"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/webinfra"
	"github.com/go-chi/chi/v5"
)

type debugController struct {
	runtime *webinfra.ServerRuntime
}

func newDebugController(runtime *webinfra.ServerRuntime) *debugController {
	return &debugController{runtime: runtime}
}

func (c *debugController) Register(r chi.Router) {
	r.Route("/debug", func(r chi.Router) {
		r.Get("/launch-time", c.launchTime)
		r.Head("/launch-time", c.launchTime)
		r.Handle("/500", http.HandlerFunc(c.testError))
	})
}

func (c *debugController) launchTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if r.Method == http.MethodGet {
		w.Write([]byte(c.runtime.StartedAt.UTC().Format(time.RFC3339Nano)))
	}
}

func (c *debugController) testError(w http.ResponseWriter, r *http.Request) {
	olhttp.Write500(w, r, errors.New("test error"))
}
