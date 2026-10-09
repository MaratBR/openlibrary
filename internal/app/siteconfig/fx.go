package siteconfig

import "go.uber.org/fx"

var FXModule = fx.Module("siteconfig",
	fx.Provide(NewStore),
	fx.Invoke(func(lifecycle fx.Lifecycle, settings Store) {
		lifecycle.Append(fx.Hook{OnStart: settings.Load})
	}),
)
