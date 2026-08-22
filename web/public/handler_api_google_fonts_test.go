package public

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGoogleFontsHandler(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll(filepath.Join(googleFontsDir, "Tangerine", "files"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Tangerine", "files", "Tangerine.ttf"), []byte("font"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Tangerine", "include.css"), []byte("tangerine-css"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Tangerine", "metadata.json"), []byte(`{"name":"Tangerine","source":"google","license":"SIL Open Font License 1.1","includeCss":"/_api/fonts/google/include?name=Tangerine","externalLink":"https://fonts.google.com/specimen/Tangerine"}`), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(googleFontsDir, "Roboto", "files"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Roboto", "include.css"), []byte("roboto-css"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Roboto", "metadata.json"), []byte(`{"name":"Roboto","source":"google","license":"Apache License 2.0","includeCss":"/_api/fonts/google/include?name=Roboto","externalLink":"https://fonts.google.com/specimen/Roboto"}`), 0o644))

	h := newAPIControllerGoogleFonts()
	for _, tc := range []struct {
		query  string
		status int
		body   string
	}{
		{"?name=Tangerine&file=Tangerine.ttf", http.StatusOK, "font"},
		{"?name=Missing&file=Tangerine.ttf", http.StatusNotFound, ""},
		{"?name=Tangerine&file=missing.ttf", http.StatusNotFound, ""},
		{"?name=Tangerine&file=../include.css", http.StatusNotFound, ""},
		{"?name=..&file=openlibrary.toml", http.StatusNotFound, ""},
		{"?name=Tangerine/files&file=Tangerine.ttf", http.StatusNotFound, ""},
	} {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/_api/fonts/google"+tc.query, nil)
			w := httptest.NewRecorder()
			h.get(w, req)
			require.Equal(t, tc.status, w.Code)
			if tc.body != "" {
				require.Equal(t, tc.body, w.Body.String())
			}
		})
	}
}

func TestGoogleFontsListHandler(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll(filepath.Join(googleFontsDir, "Roboto"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(googleFontsDir, "Broken"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Roboto", "metadata.json"), []byte(`{"name":"Roboto","source":"google","license":"Apache License 2.0","includeCss":"/_api/fonts/google/include?name=Roboto","externalLink":"https://fonts.google.com/specimen/Roboto"}`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Broken", "metadata.json"), []byte(`not-json`), 0o644))

	req := httptest.NewRequest(http.MethodGet, "/_api/fonts/google", nil)
	w := httptest.NewRecorder()
	newAPIControllerGoogleFonts().get(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.JSONEq(t, `[{"name":"Roboto","source":"google","license":"Apache License 2.0","includeCss":"/_api/fonts/google/include?name=Roboto","externalLink":"https://fonts.google.com/specimen/Roboto"}]`, w.Body.String())
}

func TestGoogleFontsIncludeHandler(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll(filepath.Join(googleFontsDir, "Tangerine"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(googleFontsDir, "Roboto"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Tangerine", "include.css"), []byte("tangerine-css"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(googleFontsDir, "Roboto", "include.css"), []byte("roboto-css"), 0o644))

	h := newAPIControllerGoogleFonts()
	for _, tc := range []struct {
		query  string
		status int
		body   string
	}{
		{"?name=Tangerine,Roboto", http.StatusOK, "tangerine-css\nroboto-css"},
		{"?name=Roboto,Tangerine", http.StatusOK, "roboto-css\ntangerine-css"},
		{"?name=Tangerine,Missing", http.StatusNotFound, "404 page not found\n"},
		{"?name=Tangerine,..", http.StatusNotFound, "404 page not found\n"},
		{"", http.StatusNotFound, "404 page not found\n"},
	} {
		t.Run(tc.query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/_api/fonts/google/include"+tc.query, nil)
			w := httptest.NewRecorder()
			h.include(w, req)
			result := w.Result()
			require.Equal(t, tc.status, result.StatusCode)
			require.Equal(t, tc.body, w.Body.String())
			if tc.status == http.StatusOK {
				require.Equal(t, "text/css; charset=utf-8", result.Header.Get("Content-Type"))
			}
		})
	}
}

func TestGoogleFontsFileAndCSSConditionalCaching(t *testing.T) {
	t.Chdir(t.TempDir())
	filesDir := filepath.Join(googleFontsDir, "Tangerine", "files")
	require.NoError(t, os.MkdirAll(filesDir, 0o755))
	fontPath := filepath.Join(filesDir, "Tangerine.ttf")
	cssPath := filepath.Join(googleFontsDir, "Tangerine", "include.css")
	require.NoError(t, os.WriteFile(fontPath, []byte("font"), 0o644))
	require.NoError(t, os.WriteFile(cssPath, []byte("css"), 0o644))
	modified := time.Date(2026, time.August, 22, 10, 0, 0, 0, time.UTC)
	require.NoError(t, os.Chtimes(fontPath, modified, modified))
	require.NoError(t, os.Chtimes(cssPath, modified, modified))

	h := newAPIControllerGoogleFonts()
	for _, tc := range []struct {
		path    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"/_api/fonts/google?name=Tangerine&file=Tangerine.ttf", h.get},
		{"/_api/fonts/google/include?name=Tangerine", h.include},
	} {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			req.Header.Set("If-Modified-Since", modified.Format(http.TimeFormat))
			w := httptest.NewRecorder()
			tc.handler(w, req)
			require.Equal(t, http.StatusNotModified, w.Code)
			require.Empty(t, w.Body.String())
			require.Equal(t, googleFontsCacheControl, w.Header().Get("Cache-Control"))
			require.Equal(t, modified.Format(http.TimeFormat), w.Header().Get("Last-Modified"))
		})
	}
}
