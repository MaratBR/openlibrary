package content

import "go.uber.org/fx"

func NewDefaultEngine() *MarkupEngine {
	return NewEngine(MarkupEngineOptions{
		TagExapanders: map[string]ExpandTag{
			"ol-widget": NewWidgetRegistry(),
		},
	})
}

var FXModule = fx.Module("content_util", fx.Provide(NewDefaultEngine))
