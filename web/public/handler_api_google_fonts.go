package public

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const googleFontsDir = "oldata/fonts/google"

const googleFontsCacheControl = "public, max-age=86400"

type apiControllerGoogleFonts struct{}

type googleFontMetadata struct {
	Name         string `json:"name"`
	Source       string `json:"source"`
	License      string `json:"license"`
	IncludeCSS   string `json:"includeCss"`
	ExternalLink string `json:"externalLink"`
}

func newAPIControllerGoogleFonts() *apiControllerGoogleFonts { return &apiControllerGoogleFonts{} }

func (c *apiControllerGoogleFonts) Register(r chi.Router) {
	r.Get("/fonts/google", c.get)
	r.Get("/fonts/google/include", c.include)
}

func (c *apiControllerGoogleFonts) include(w http.ResponseWriter, r *http.Request) {
	names := strings.Split(r.URL.Query().Get("name"), ",")
	if len(names) == 0 || len(names) > 100 {
		http.NotFound(w, r)
		return
	}
	for _, name := range names {
		if !safeFontPathComponent(name) {
			http.NotFound(w, r)
			return
		}
	}

	root, err := os.OpenRoot(googleFontsDir)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer root.Close()

	stylesheets := make([][]byte, 0, len(names))
	var lastModified time.Time
	for _, name := range names {
		f, err := root.Open(filepath.ToSlash(filepath.Join(name, "include.css")))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		stat, statErr := f.Stat()
		contents, readErr := io.ReadAll(io.LimitReader(f, (1<<20)+1))
		closeErr := f.Close()
		if statErr != nil || !stat.Mode().IsRegular() || readErr != nil || closeErr != nil || len(contents) > 1<<20 {
			http.NotFound(w, r)
			return
		}
		if stat.ModTime().After(lastModified) {
			lastModified = stat.ModTime()
		}
		stylesheets = append(stylesheets, contents)
	}

	var combined bytes.Buffer
	for i, stylesheet := range stylesheets {
		if i > 0 {
			combined.WriteByte('\n')
		}
		combined.Write(stylesheet)
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", googleFontsCacheControl)
	http.ServeContent(w, r, "include.css", lastModified, bytes.NewReader(combined.Bytes()))
}

func (c *apiControllerGoogleFonts) get(w http.ResponseWriter, r *http.Request) {
	name, filename := r.URL.Query().Get("name"), r.URL.Query().Get("file")
	if name == "" && filename == "" {
		c.list(w, r)
		return
	}
	if !safeFontPathComponent(name) || !safeFontPathComponent(filename) {
		http.NotFound(w, r)
		return
	}

	root, err := os.OpenRoot(googleFontsDir)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer root.Close()
	f, err := root.Open(filepath.ToSlash(filepath.Join(name, "files", filename)))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", googleFontsCacheControl)
	http.ServeContent(w, r, filename, stat.ModTime(), f)
}

func (c *apiControllerGoogleFonts) list(w http.ResponseWriter, r *http.Request) {
	root, err := os.OpenRoot(googleFontsDir)
	if err != nil {
		writeGoogleFontsJSON(w, []googleFontMetadata{})
		return
	}
	defer root.Close()
	dir, err := root.Open(".")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	entries, err := dir.ReadDir(-1)
	closeErr := dir.Close()
	if closeErr != nil {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fonts := make([]googleFontMetadata, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || !safeFontPathComponent(entry.Name()) {
			continue
		}
		f, err := root.Open(filepath.ToSlash(filepath.Join(entry.Name(), "metadata.json")))
		if err != nil {
			continue
		}
		var metadata googleFontMetadata
		decodeErr := json.NewDecoder(io.LimitReader(f, 64<<10)).Decode(&metadata)
		closeErr := f.Close()
		if decodeErr != nil || closeErr != nil || metadata.Name != entry.Name() {
			continue
		}
		fonts = append(fonts, metadata)
	}
	writeGoogleFontsJSON(w, fonts)
}

func writeGoogleFontsJSON(w http.ResponseWriter, fonts []googleFontMetadata) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_ = json.NewEncoder(w).Encode(fonts)
}

func safeFontPathComponent(value string) bool {
	return value != "" && value != "." && value != ".." && filepath.Base(value) == value && !strings.ContainsAny(value, "/\\\x00\r\n")
}
