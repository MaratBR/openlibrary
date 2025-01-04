package uiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
)

type dataKeyType struct{}

var dataKey dataKeyType

type Handler struct {
	devFrontEndServerProxy http.Handler
	router                 chi.Router
	defaultDataLoader      func(r *http.Request, data *Data)
}

func (c *Handler) PreloadData(pattern string, f func(r *http.Request, data *Data)) {
	c.router.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(dataKey).(*Data)
		if c.defaultDataLoader != nil {
			c.defaultDataLoader(r, data)
		}
		f(r, data)
		r = r.WithContext(context.WithValue(r.Context(), dataKey, data))
		c.devFrontEndServerProxy.ServeHTTP(w, r)
	})
}

func (c *Handler) DefaultPreloadData(f func(r *http.Request, data *Data)) {
	c.defaultDataLoader = f
}

func (c *Handler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if c.defaultDataLoader == nil {
		c.devFrontEndServerProxy.ServeHTTP(w, r)
	} else {
		data := newData()
		c.defaultDataLoader(r, data)
		r = r.WithContext(context.WithValue(r.Context(), dataKey, data))
		c.devFrontEndServerProxy.ServeHTTP(w, r)
	}
}

func (c *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}

func NewDevStaticHandler(
	config *koanf.Koanf,
) *Handler {
	address := fmt.Sprintf("%s://%s:%d", config.String("frontend-proxy.target-protocol"), config.String("frontend-proxy.target-host"), config.Int("frontend-proxy.target-port"))
	c := &Handler{
		router: chi.NewRouter(),
		devFrontEndServerProxy: newProxy(address, DevServerOptions{
			GetInjectedHTMLSegment: func(r *http.Request) []byte {
				data := r.Context().Value(dataKey)
				if data == nil {
					return nil
				}

				return getInjectedHTMLSegment(data.(*Data))
			},
		}),
	}

	c.router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := newData()
			r = r.WithContext(context.WithValue(r.Context(), dataKey, data))
			h.ServeHTTP(w, r)
		})
	})

	c.router.NotFound(c.notFoundHandler)

	return c
}
