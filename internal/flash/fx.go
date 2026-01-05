package flash

import "go.uber.org/fx"

var FXModule = fx.Module("flash", fx.Provide(
	MakeMiddleware,
))
