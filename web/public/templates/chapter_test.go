package templates

import (
	"bytes"
	"context"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app/content"
	"github.com/stretchr/testify/require"
)

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
