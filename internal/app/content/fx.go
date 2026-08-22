package content

import "go.uber.org/fx"

func NewDefaultEngine() *MarkupEngine {
	return NewEngine(MarkupEngineOptions{
		TagExapanders: map[string]ExpandTag{
			"ol-widget": NewWidgetRegistry(),
		},
		TextColorFilter:     func(color string) bool { return true },
		URLFilter:           func(url string) bool { return true },
		AllowedFontFamilies: []string{},
		AllowedFontSizes: []string{
			"1em",
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
		},
	})
}

var FXModule = fx.Module("content_util", fx.Provide(NewDefaultEngine))
