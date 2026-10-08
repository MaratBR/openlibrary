package templates

import (
	"bytes"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

func TestCensoredBookPreviewIsInitiallyBlurredAndInert(t *testing.T) {
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../../translations/en.toml"}}, zap.NewNop().Sugar())
	for _, mode := range []app.CensorMode{"none", "hide", "censor"} {
		ctx := app.WithContentPreferences(localizedContext(provider), &app.ContentPreferences{Mode: mode, TagIDs: []int64{1}})
		app.ApplyBookTagPreferences(ctx, 10, []int64{1})
		var out bytes.Buffer
		require.NoError(t, BookPreviewPartial(app.BookDetailsDto{ID: 10, Name: "Book", Summary: "summary"}).Render(ctx, &out))
		html := out.String()
		if mode != "none" {
			require.Contains(t, html, `data-censored-book="10"`)
			require.Contains(t, html, `filter:blur(12px)`)
			require.Contains(t, html, `inert aria-hidden="true"`)
			require.Contains(t, html, `@click="reveal()"`)
		} else {
			require.NotContains(t, html, "data-censored-book")
		}
	}
}
