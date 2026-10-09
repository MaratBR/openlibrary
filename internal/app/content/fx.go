package content

import (
	"strings"

	"go.uber.org/fx"
)

func NewDefaultEngine() *MarkupEngine {
	return NewEngine(MarkupEngineOptions{
		TagExapanders: map[string]ExpandTag{
			"ol-widget": NewWidgetRegistry(),
		},
		TextColorFilter: func(color string) bool { return true },
		URLFilter:       func(url string) bool { return true },
		FontFamilyFilter: func(fontFamily string) bool {
			return strings.TrimSpace(fontFamily) != ""
		},
		AllowedFontSizes: DefaultAllowedFontSizes(),
	})
}

// DefaultAllowedFontSizes is shared by chapter sanitization and the editor UI.
func DefaultAllowedFontSizes() []string {
	return []string{
		"8px", "10px", "12px", "14px", "16px", "18px", "20px",
		"24px", "28px", "32px", "36px", "48px", "72px",
		"1em",
		"1.25em",
		"1.25em",
		"1.6em",
		"1.8em",
		"2em",
		"3em",
		"4em",
		"5em",
		"6em",
		"7em",
		"8em",
		"9em",
		"10em",
		"11em",
		"12em",
	}
}

var FXModule = fx.Module("content_util", fx.Provide(NewDefaultEngine))
