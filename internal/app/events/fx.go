package events

import "go.uber.org/fx"

var FXModule = fx.Module("app_events", fx.Provide(NewHub))
