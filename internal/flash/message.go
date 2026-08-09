package flash

import (
	"context"
	"io"

	"github.com/MaratBR/openlibrary/internal/app/content"
	"github.com/a-h/templ"
)

type Message struct {
	Text string `json:"text"`
}

// Render implements Message.
func (t Message) Render(ctx context.Context, w io.Writer) error {
	return templ.Raw(content.SanitizeHtml(t.Text)).Render(ctx, w)
}

func Text(text string) Message {
	return Message{Text: text}
}
