package olhttp

import (
	"context"
	"go.uber.org/zap"
	"net/http"

	"github.com/a-h/templ"
)

func WriteTemplate(w http.ResponseWriter, ctx context.Context, t templ.Component) {
	err := t.Render(ctx, w)
	if err != nil {
		zap.S().Errorw("error while writing to the client", "err", err)
		return
	}
	//	w.Write([]byte(`
	//
	// <!--
	//
	//	 OOOOO  L
	//	O     O L
	//	O     O L
	//	O     O L
	//	 OOOOO  LLLLLLLL
	//
	// -->`))
}
