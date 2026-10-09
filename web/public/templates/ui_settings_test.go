package templates

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReaderWidthCookies(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  int16
	}{
		{"48", 48}, {"84", 84}, {"100", 100}, {"49", 72}, {"102", 72}, {"65584", 72}, {"bad", 72},
	} {
		t.Run(tc.value, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.AddCookie(&http.Cookie{Name: "reader_content_width", Value: tc.value})
			r.AddCookie(&http.Cookie{Name: "reader_font_family", Value: "sans"})
			got := GetReaderPreferencesFromCookies(r)
			if got.ContentWidth != tc.want || got.FontFamily != "sans" {
				t.Fatalf("unexpected preferences: %+v", got)
			}
		})
	}
}

func TestOLEDReaderThemeCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "reader_theme", Value: "oled"})
	if got := GetReaderPreferencesFromCookies(r); got.Theme != "oled" {
		t.Fatalf("OLED preference lost: %+v", got)
	}
}
