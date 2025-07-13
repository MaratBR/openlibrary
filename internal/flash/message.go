package flash

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

type Message struct {
	Text string `json:"text"`
}

// Render implements Message.
func (t Message) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(templ.EscapeString(string(t.Text))))
	return err
}

func Text(text string) Message {
	return Message{Text: text}
}
