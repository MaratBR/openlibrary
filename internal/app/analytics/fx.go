package analytics

import "go.uber.org/fx"

var FXModule = fx.Module(
	"ol_app_analytics",

	eventModule,
	atomicModule,
	metricModule,
	popularityModule,

	fx.Provide(
		NewAnalyticsCounters,
	),
)
