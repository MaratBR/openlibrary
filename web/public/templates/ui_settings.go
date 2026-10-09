package templates

import (
	"net/http"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app"
)

func GetReaderPreferencesFromCookies(r *http.Request) app.ReaderPreferences {
	preferences := app.DefaultReaderPreferences()

	fontSizeCookie, fontSizeCookieErr := r.Cookie("reader_font_size")
	if fontSizeCookieErr != nil {
		// Preserve the font-size preference used by the previous reader UI.
		fontSizeCookie, fontSizeCookieErr = r.Cookie("ifs")
	}
	if fontSizeCookieErr == nil {
		if value, parseErr := strconv.Atoi(fontSizeCookie.Value); parseErr == nil && value >= 12 && value <= 48 {
			preferences.FontSize = int16(value)
		}
	}
	if cookie, err := r.Cookie("reader_font_family"); err == nil {
		preferences.FontFamily = cookie.Value
	}
	if cookie, err := r.Cookie("reader_page_color"); err == nil {
		preferences.PageColor = cookie.Value
	}
	if cookie, err := r.Cookie("reader_content_width"); err == nil {
		if value, parseErr := strconv.Atoi(cookie.Value); parseErr == nil && value >= 48 && value <= 100 && value%2 == 0 {
			preferences.ContentWidth = int16(value)
		}
	}
	if cookie, err := r.Cookie("reader_theme"); err == nil {
		preferences.Theme = cookie.Value
	}

	if err := preferences.Validate(); err != nil {
		return app.DefaultReaderPreferences()
	}
	return preferences
}
