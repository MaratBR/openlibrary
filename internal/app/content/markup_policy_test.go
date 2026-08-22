package content

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestMarkupEngineCleanAllowsTextAlignmentValues(t *testing.T) {
	t.Parallel()

	for _, alignment := range []string{"left", "right", "center", "justify"} {
		t.Run(alignment, func(t *testing.T) {
			t.Parallel()

			input := fmt.Sprintf(`<p style="text-align: %s">aligned</p>`, alignment)
			assertSanitizedMarkup(t, NewEngine(MarkupEngineOptions{}), input, input)
		})
	}
}

func TestMarkupEngineCleanAllowsWidgetsWithOptionalText(t *testing.T) {
	t.Parallel()

	engine := NewDefaultEngine()
	for _, test := range []struct {
		input string
		want  string
	}{
		{`<ol-widget name="example">TEXT_HERE</ol-widget>`, `<ol-widget name="example">TEXT_HERE</ol-widget>`},
		{`<ol-widget name="example"></ol-widget>`, `<ol-widget name="example"></ol-widget>`},
		{`<ol-widget name="example"/>`, `<ol-widget name="example"></ol-widget>`},
	} {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			assertSanitizedMarkup(t, engine, test.input, test.want)
		})
	}
}

func TestMarkupEngineCleanAllowsInlineFormatting(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{
		TextColorFilter: func(color string) bool { return color == "#123456" },
	})
	assertSanitizedMarkup(t, engine,
		`<b>bold</b><i>italic</i><s>strikethrough</s><sup>superscript</sup><sub>subscript</sub><span style="color: #123456">color</span><ruby>漢<rt>かん</rt></ruby><mark>mark</mark><code>code</code>`,
		`<b>bold</b><i>italic</i><s>strikethrough</s><sup>superscript</sup><sub>subscript</sub><span style="color: #123456">color</span><ruby>漢<rt>かん</rt></ruby><mark>mark</mark><code>code</code>`,
	)
}

func TestMarkupEngineCleanFiltersTextColor(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{
		TextColorFilter: func(color string) bool { return color == "#123456" },
	})
	assertSanitizedMarkup(t, engine,
		`<span style="color: #654321">color</span>`,
		`<span>color</span>`,
	)
}

func TestMarkupEngineCleanAllowsConfiguredFontValuesOnly(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{
		AllowedFontFamilies: []string{"Arial"},
		AllowedFontSizes:    []string{"16px"},
	})
	assertSanitizedMarkup(t, engine,
		`<span style="font-family: Arial; font-size: 16px">styled</span>`,
		`<span style="font-family: Arial; font-size: 16px">styled</span>`,
	)
	assertSanitizedMarkup(t, engine,
		`<span style="font-family: Unapproved; font-size: 17px">styled</span>`,
		`<span>styled</span>`,
	)
}

func TestMarkupEngineCleanReturnsSeenFonts(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{
		AllowedFontFamilies: []string{"Arial", "Georgia", "Arial, Georgia"},
	})
	got, err := engine.Clean(`<p style="font-family: Arial">one</p><span style="font-family: Arial, Georgia">two</span><em style="font-family: Arial">three</em><b style="font-family: Unapproved">four</b>`)
	if err != nil {
		t.Fatalf("Clean() error = %v", err)
	}
	if want := []string{"Arial", "Georgia"}; !slices.Equal(got.Fonts, want) {
		t.Errorf("Clean().Fonts = %#v, want %#v", got.Fonts, want)
	}
}

func TestMarkupEngineCleanFiltersLinkURLs(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{
		URLFilter: func(url string) bool {
			return strings.HasPrefix(url, "https://allowed.example/") || strings.HasPrefix(url, "/")
		},
	})
	for _, test := range []struct {
		input string
		want  string
	}{
		{`<a href="https://allowed.example/page">absolute</a>`, `<a href="https://allowed.example/page" rel="nofollow">absolute</a>`},
		{`<a href="https://allowed.example/page" target="_blank">new window</a>`, `<a href="https://allowed.example/page" rel="nofollow noopener" target="_blank">new window</a>`},
		{`<a href="/relative/page">relative</a>`, `<a href="/relative/page" rel="nofollow">relative</a>`},
	} {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			assertSanitizedMarkup(t, engine, test.input, test.want)
		})
	}
	assertSanitizedMarkup(t, engine, `<a href="https://blocked.example/page">blocked</a>`, `<a>blocked</a>`)
}

func TestMarkupEngineCleanFiltersImageURLs(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{
		URLFilter: func(url string) bool { return strings.HasPrefix(url, "https://images.example/") },
	})
	assertSanitizedMarkup(t, engine,
		`<img src="https://images.example/cover.png" alt="Cover"/>`,
		`<img src="https://images.example/cover.png" alt="Cover"/>`,
	)
	assertSanitizedMarkup(t, engine,
		`<img src="https://blocked.example/cover.png" alt="Cover">`,
		`<img alt="Cover"/>`,
	)
}

func TestMarkupEngineCleanAllowsCodeBlocksAndTables(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{})
	assertSanitizedMarkup(t, engine, `<pre>plain code</pre>`, `<pre>plain code</pre>`)
	assertSanitizedMarkup(t, engine, `<pre data-lang="go">fmt.Println()</pre>`, `<pre data-lang="go">fmt.Println()</pre>`)
	assertSanitizedMarkup(t, engine,
		`<table><thead><tr><th>Header</th></tr></thead><tbody><tr><td>Cell</td></tr></tbody></table>`,
		`<table><thead><tr><th>Header</th></tr></thead><tbody><tr><td>Cell</td></tr></tbody></table>`,
	)
}

func assertSanitizedMarkup(t *testing.T, engine *MarkupEngine, input, want string) {
	t.Helper()

	got, err := engine.Clean(input)
	if err != nil {
		t.Fatalf("Clean(%q) error = %v", input, err)
	}
	if got.Sanitized != want {
		t.Errorf("Clean(%q) = %q, want %q", input, got.Sanitized, want)
	}
}
