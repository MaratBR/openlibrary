package frontend

import (
	"embed"
	"io/fs"
	"net/http"
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
	// Set cache headers assuming the response does not change and was created at assetsTime
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("Last-Modified", assetsTime.UTC().Format(http.TimeFormat))

	a.fileServer.ServeHTTP(w, r)
}

type embedAssetsFS struct {
	inner fs.FS
}

// Open implements fs.FS.
func (e embedAssetsFS) Open(name string) (fs.File, error) {
	name = "embed-assets/" + name
	file, err := e.inner.Open(name)
	return file, err
}

func EmbedAssetsFS() fs.FS {
	return embedAssetsFS{inner: folder}
}

func EmbedAssets(fs fs.FS) http.Handler {
	return assetsHandler{
		fileServer: http.FileServerFS(fs),
	}
}
