package tags

import "go.uber.org/fx"

var FXModule = fx.Module("tags",
	fx.Provide(NewTagsService),
	fx.Invoke(startTagsPopulateService),
)
