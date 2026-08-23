package bookfont

import "go.uber.org/fx"

var FXModule = fx.Module("bookfont", fx.Provide(NewPolicy))
