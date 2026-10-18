package templates

import (
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

func TestSearchTitleEscapesHighlightedText(t *testing.T) {
	provider := i18n.NewLocaleProvider(language.English, false, map[language.Tag][]string{language.English: {"../../../translations/en.toml"}}, zap.NewNop().Sugar())
	var output strings.Builder
	title := `<script>alert(1)</script> & Book`
	err := bookSearchCard(1, title, "book", "", app.BookCover{}, "Author", [16]byte{}, nil, []app.SearchTextPart{
		{Text: `<script>alert(1)</script>`, Matched: true}, {Text: " & Book"},
	}).Render(localizedContext(provider), &output)
	if err != nil {
		t.Fatal(err)
	}
	rendered := output.String()
	if !strings.Contains(rendered, "<mark>&lt;script&gt;alert(1)&lt;/script&gt;</mark>") || strings.Contains(rendered, "<script>alert(1)</script>") {
		t.Fatalf("unsafe or missing highlight: %s", rendered)
	}
}
