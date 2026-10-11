package userdata

import "go.uber.org/fx"

var FXModule = fx.Module("userdata", fx.Provide(NewUserDataService))
