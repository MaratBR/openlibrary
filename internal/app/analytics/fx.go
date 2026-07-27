package analytics

import "go.uber.org/fx"

var FXModule = fx.Module(
	"ol_app_analytics",

	fx.Provide(
		newEventBackgroundService,
		newEventRepository,
		fx.Private,
	),

	fx.Provide(
		NewAnalyticsCounters,
		func() ViewsService {
			return &analyticsViewsDummyService{}
		},

		NewAtomic,
		newDedupedEventSink,
	),

	fx.Invoke(
		func(*eventBackgroundService) {},
	),
)
