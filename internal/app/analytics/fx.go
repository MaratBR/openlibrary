package analytics

import "go.uber.org/fx"

var FXModule = fx.Module(
	"ol_app_analytics",

	eventModule,
	atomicModule,
	metricModule,

	fx.Provide(
		NewAnalyticsCounters,
		func() ViewsService {
			return &analyticsViewsDummyService{}
		},
	),
)
