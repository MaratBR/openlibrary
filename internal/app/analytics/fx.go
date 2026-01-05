package analytics

import "go.uber.org/fx"

var FXModule = fx.Module("ol_app_analytics", fx.Provide(
	NewAnalyticsCounters,
	NewAnalyticsViewsService,
	NewAnalyticsBackgroundService,
), fx.Invoke(func(srv *AnalyticsBackgroundService) {}))
