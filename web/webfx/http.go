package webfx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type MountableHandler interface {
	http.Handler

	MountAt() string
}

type WebController interface {
	Register(r chi.Router)
}

func AsController(fn any) any {
	return fx.Annotate(
		fn,
		fx.As(new(WebController)),
		fx.ResultTags(`group:"routes"`),
	)
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
