package olhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDocInitialTheme(t *testing.T) {
	for _, theme := range []string{"", "system", "light", "dark"} {
		t.Run("theme="+theme, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if theme != "" {
				r.AddCookie(&http.Cookie{Name: "_theme", Value: theme})
			}
			ctx := context.WithValue(r.Context(), contextKeyRequest, r)
			var output strings.Builder
			if err := Doc(DocProps{CSS: []string{"/test.css"}}).Render(ctx, &output); err != nil {
				t.Fatal(err)
			}
			html := output.String()
			systemScript := strings.Index(html, "window.matchMedia('(prefers-color-scheme: dark)')")
			if theme == "" || theme == "system" {
				if systemScript < 0 || systemScript > strings.Index(html, `<link rel="stylesheet"`) {
					t.Fatal("system theme must resolve before stylesheets load")
				}
			} else if systemScript >= 0 {
				t.Fatal("explicit theme must not be overridden by system preference")
			}
			if strings.Contains(html, `class="theme-default dark"`) != (theme == "dark") {
				t.Fatal("incorrect server-rendered theme class")
			}
		})
	}
}

func TestDocReaderThemeBeforeStylesheets(t *testing.T) {
	for _, readerTheme := range []string{"dark", "oled", "light"} {
		t.Run(readerTheme, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			websiteTheme := "light"
			if readerTheme == "light" {
				websiteTheme = "dark"
			}
			r.AddCookie(&http.Cookie{Name: "_theme", Value: websiteTheme})
			ctx := context.WithValue(r.Context(), contextKeyRequest, r)
			var output strings.Builder
			if err := Doc(DocProps{ReaderTheme: readerTheme, CSS: []string{"/test.css"}}).Render(ctx, &output); err != nil {
				t.Fatal(err)
			}
			html := output.String()
			dark := readerTheme != "light"
			if strings.Contains(html, `class="theme-default dark"`) != dark {
				t.Fatal("incorrect initial reader dark class")
			}
			expectedTheme := "light"
			if dark {
				expectedTheme = "dark"
			}
			for _, attribute := range []string{`data-theme="` + expectedTheme + `"`, `data-reader-theme="` + readerTheme + `"`} {
				index := strings.Index(html, attribute)
				if index < 0 || index > strings.Index(html, "<link") {
					t.Fatalf("missing early attribute %s", attribute)
				}
			}
			if strings.Contains(html, "window.matchMedia") {
				t.Fatal("explicit reader theme must not be overridden by system preference")
			}
		})
	}
}
