package reader

import "go.uber.org/fx"

var FXModule = fx.Module("reader", fx.Provide(
	NewReaderPreferencesService,
	NewReaderFontService,
))
