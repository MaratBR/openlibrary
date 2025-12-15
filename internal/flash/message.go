package flash

import (
	"context"
	"io"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/a-h/templ"
)

type Message struct {
	Text string `json:"text"`
}

// Render implements Message.
func (t Message) Render(ctx context.Context, w io.Writer) error {
	return templ.Raw(app.SanitizeHtml(t.Text)).Render(ctx, w)
}

func Text(text string) Message {
	return Message{Text: text}
}
