package notifications

import "go.uber.org/fx"

var FXModule = fx.Module("app_notifications",
	fx.Provide(
		NewSQLPoller,
		NewSQLStore,
		NewSimpleSink,
	))
