package webinfra

import (
	"net/http"

	"go.uber.org/fx"
)

type MountableHandler interface {
	http.Handler

	MountAt() string
}

func AsMountableHandler(fn any) any {
	return fx.Annotate(
		fn,
		fx.As(new(MountableHandler)),
		fx.ResultTags(`group:"root_handlers"`),
	)
}

func ProviderMountables() fx.Annotation {
	return fx.ParamTags(`group:"root_handlers"`)
}
