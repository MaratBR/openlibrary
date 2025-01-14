package csrf

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
)

func CSRFInputString(ctx context.Context) string {
	v := ctx.Value(csrfTokenKey)
	if v == nil {
		return ""
	}
	if token, ok := v.(string); ok {
		return fmt.Sprintf("<input type=\"hidden\" name=\"__csrf\" value=\"%s\"/>", token)
	} else {
		return ""
	}
}

func CSRFInputTempl(ctx context.Context) templ.Component {
	str := CSRFInputString(ctx)
	if str == "" {
		return templ.NopComponent
	} else {
		return templ.Raw(str)
	}
}
