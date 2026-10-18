package notifications

import (
	"github.com/MaratBR/openlibrary/internal/app/events"
	"go.uber.org/fx"
)

var FXModule = fx.Module("app_notifications",
	fx.Provide(
		NewSQLPoller,
		NewStreamReader,
		func(reader *StreamReader) Service { return reader },
		fx.Annotate(func(reader *StreamReader) events.Source { return reader }, fx.ResultTags(`group:"user_sse_sources"`)),
		NewSQLStore,
		NewSimpleSink,
	))
