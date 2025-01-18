package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

type htmlTemplateElement struct {
	Data any
	ID   string
}

// Render implements templ.Component.
func (h *htmlTemplateElement) Render(ctx context.Context, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "<template type=\"application/json\" id=\"%s\">", templ.EscapeString(h.ID)); err != nil {
		return err
	}

	{
		enc := json.NewEncoder(w)
		if err := enc.Encode(h.Data); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("</template>")); err != nil {
		return err
	}

	return nil
}

func jsonTemplate(id string, data any) templ.Component {
	return &htmlTemplateElement{
		ID:   id,
		Data: data,
	}
}
