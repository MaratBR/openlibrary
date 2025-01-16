package publicui

import (
	"embed"
	"net/http"
	"net/url"
	"time"
)

//go:embed embed-assets/*
var folder embed.FS
var assetsTime = time.Now()

type assetsHandler struct {
	fileServer http.Handler
}

// ServeHTTP implements http.Handler.
func (a assetsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL

	r2.URL.Path = "embed-assets/" + r2.URL.Path
	r2.URL.RawPath = "/embed-assets/" + r2.URL.RawPath

	// Set cache headers assuming the response does not change and was created at assetsTime
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Last-Modified", assetsTime.UTC().Format(http.TimeFormat))

	a.fileServer.ServeHTTP(w, r2)
}

func newAssetsHandler() http.Handler {
	return assetsHandler{
		fileServer: http.FileServerFS(folder),
	}
}
