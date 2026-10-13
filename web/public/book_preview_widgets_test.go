package public

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/content"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type previewWidgetBooks struct {
	app.BookService
	book    app.BookDetailsDto
	err     error
	queries []app.GetBookQuery
}

func (s *previewWidgetBooks) GetBookDetails(_ context.Context, query app.GetBookQuery) (app.BookDetailsDto, error) {
	s.queries = append(s.queries, query)
	return s.book, s.err
}

func renderPreviewWidgetFixture(t *testing.T, service *previewWidgetBooks, markup string) (string, map[string]string) {
	t.Helper()
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../translations/en.toml"}}, zap.NewNop().Sugar())
	router := chi.NewRouter()
	router.Use(provider.Middleware)
	var result string
	var previews map[string]string
	router.Get("/", func(_ http.ResponseWriter, r *http.Request) {
		expanded, err := content.NewDefaultEngine().Expand(markup)
		require.NoError(t, err)
		result, previews, err = renderBookPreviewWidgets(r, expanded, service)
		require.NoError(t, err)
	})
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	return result, previews
}

func TestBookPreviewWidgetsRenderServerSide(t *testing.T) {
	service := &previewWidgetBooks{book: app.BookDetailsDto{
		ID: 42, Name: "A compact & complete book", Summary: "<p>Server-rendered summary.</p>",
		Author: app.BookDetailsAuthorDto{Name: "An author"},
		Cover:  app.BookCover{URL: "/covers/book.webp"},
		Tags:   []app.DefinedTagDto{{ID: 7, Name: "Fiction"}},
	}}
	widget := `<ol-widget name="book-preview" data='{"bookId":"42"}'><p>Forged data</p></ol-widget>`
	result, previews := renderPreviewWidgetFixture(t, service, "<p>Before</p>"+widget+widget+"<p>After</p>")
	for _, expected := range []string{"A compact &amp; complete book", "An author", "Server-rendered summary", "/covers/book.webp", "Fiction", `href="/book/42"`, "Open book", "<p>Before</p>", "<p>After</p>"} {
		require.Contains(t, result, expected)
	}
	require.NotContains(t, result, "Forged data")
	require.NotContains(t, result, "Loading book preview")
	require.Contains(t, previews["42"], `class="BookCardPreview"`)
	require.Len(t, service.queries, 1, "duplicate widgets share a request-scoped lookup")
	require.Equal(t, int64(42), service.queries[0].ID)
	require.False(t, service.queries[0].ActorUserID.Valid)
	require.False(t, service.queries[0].AdminOverride)
	if path := os.Getenv("BOOK_PREVIEW_FIXTURE"); path != "" {
		require.NoError(t, os.WriteFile(path, []byte(result), 0600))
	}
}

func TestBookPreviewWidgetsDoNotExposeUnavailableBooks(t *testing.T) {
	for _, err := range []error{app.ErrTypeBookNotFound.New("missing"), app.ErrTypeBookPrivated.New("private")} {
		service := &previewWidgetBooks{err: err, book: app.BookDetailsDto{Name: "Private title", Summary: "Private summary"}}
		result, previews := renderPreviewWidgetFixture(t, service, `<ol-widget name="book-preview" data='{"bookId":"42"}'></ol-widget>`)
		require.Contains(t, result, "Book preview unavailable")
		require.NotContains(t, result, "Private title")
		require.NotContains(t, result, "Private summary")
		require.Contains(t, previews["42"], "Book preview unavailable")
	}
}

func TestBookPreviewWidgetsRejectInvalidReferencesBeforeLookup(t *testing.T) {
	service := &previewWidgetBooks{}
	_, previews := renderPreviewWidgetFixture(t, service, `<ol-widget name="book-preview" data='{"bookId":"-1"}'></ol-widget>`)
	require.Empty(t, service.queries)
	require.Empty(t, previews)
}

func TestBookPreviewWidgetsDoNotRecursivelyRenderSummaryReferences(t *testing.T) {
	widget := `<ol-widget name="book-preview" data='{"bookId":"42"}'></ol-widget>`
	service := &previewWidgetBooks{book: app.BookDetailsDto{ID: 42, Name: "Book", Summary: widget}}
	_, _ = renderPreviewWidgetFixture(t, service, widget)
	require.Len(t, service.queries, 1)
}
