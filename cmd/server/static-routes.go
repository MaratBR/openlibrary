package main

import (
	"log/slog"
	"net/http"

	uiserver "github.com/MaratBR/openlibrary/cmd/server/ui-server"
)

var devFrontEndServerProxy = uiserver.NewDevServerProxy(uiserver.DevServerOptions{
	GetInjectedHTMLSegment: getInjectedHTMLSegment,
})

type staticController struct{}

func newStaticController() staticController {
	return staticController{}
}

func (staticController) DevProxyIndex(w http.ResponseWriter, r *http.Request) {
	devFrontEndServerProxy.ServeHTTP(w, r)
}

func (c staticController) PreloadData(preload func(r *http.Request, serverData *serverData) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, ok := getServerData(r)
		if ok {
			err := preload(r, data)
			if err != nil {
				slog.Error("failed to preload data", "err", err)
			}
		}

		c.DevProxyIndex(w, r)
	}
}

func (staticController) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}
