package templates

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/content"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/frontend"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

func TestEditorReadingSettings(t *testing.T) {
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../../translations/en.toml"}}, zap.NewNop().Sugar())
	preferences := app.ReaderPreferences{ContentWidth: 84, FontSize: 22, FontFamily: app.ReaderFontSans, PageColor: app.ReaderPageSurface, Theme: app.ReaderThemeDark}
	var out bytes.Buffer
	router := chi.NewRouter()
	frontend.AttachAssetsInliningHandler(os.DirFS("../../frontend/embed-assets"), "embed-assets", router)
	router.Use(provider.Middleware, olhttp.ReqCtxMiddleware)
	router.Get("/", func(_ http.ResponseWriter, r *http.Request) {
		require.NoError(t, ChapterContentIframe(preferences, app.NewReaderFontService().List()).Render(r.Context(), &out))
	})
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	for _, expected := range []string{`data-content-width="84"`, `x-model.number="contentWidth"`, `data-font-stack="Merriweather, serif"`, `data-theme-preview="dark"`, `data-font-size="22"`, `data-font-family="sans"`, `data-page-color="surface"`, `data-reader-theme="dark"`, `data-authenticated="true"`, `id="BlockEditorWrap"`, `id="ChapterContent"`, `x-model="fontFamily"`, `x-model="pageColor"`, `x-model="readerTheme"`, `x-bind="increaseFont"`, `x-bind="decreaseFont"`, `x-bind="closeButton"`} {
		require.Contains(t, out.String(), expected)
	}
	if path := os.Getenv("EDITOR_READING_FIXTURE"); path != "" {
		require.NoError(t, os.WriteFile(path, out.Bytes(), 0600))
	}
}

func TestChapterViewerPreservesSavedFontsAndSizes(t *testing.T) {
	saved, err := content.NewDefaultEngine().Clean(`<p><span style="font-family: &quot;Playfair Display&quot;; font-size: 24px">custom text</span></p>`)
	require.NoError(t, err)
	require.Equal(t, []string{"Playfair Display"}, saved.Fonts)
	require.Equal(t, []string{"/_api/fonts/google/include?name=Playfair+Display"}, chapterFontStylesheets(saved.Fonts))

	var out bytes.Buffer
	require.NoError(t, ChapterContent(saved.Sanitized).Render(context.Background(), &out))
	require.Contains(t, out.String(), saved.Sanitized)
	require.Contains(t, out.String(), "font-size: 24px")
	require.NotContains(t, out.String(), "&amp;#34;")
}
