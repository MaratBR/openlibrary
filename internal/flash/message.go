package flash

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

type Message templ.Component

type textFlash string

// Render implements Message.
func (t textFlash) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(templ.EscapeString(string(t))))
	return err
}

func Text(text string) Message {
	return textFlash(text)
}
