package frontend

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type inlineAssetsID string

func AttachAssetsInliningHandler(fs fs.FS, name string, r chi.Router) {
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := inlineAssetsID(name)
			r = r.WithContext(context.WithValue(r.Context(), key, fs))

			h.ServeHTTP(w, r)
		})
	})
}

func writeAsset(ctx context.Context, groupName, name string, w io.Writer) error {
	v := ctx.Value(inlineAssetsID(groupName))
	if v == nil {
		return errors.New("cannot find inline assets handler with name " + groupName)
	}

	fs := v.(fs.FS)

	file, err := fs.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return nil
}

type inlineCSSAsset struct {
	name      string
	groupName string
}

// Render implements templ.Component.
func (i inlineCSSAsset) Render(ctx context.Context, w io.Writer) error {
	if _, err := w.Write([]byte("<style>")); err != nil {
		return err
	}

	err := writeAsset(ctx, i.groupName, i.name, w)
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte("</style>")); err != nil {
		return err
	}

	return nil
}

func InlineCSSAsset(ctx context.Context, groupName, name string) templ.Component {
	return inlineCSSAsset{groupName: groupName, name: name}
}

func InlineJSModuleAsset(ctx context.Context, groupName, name string) templ.Component {
	return inlineCSSAsset{groupName: groupName, name: name}
}
