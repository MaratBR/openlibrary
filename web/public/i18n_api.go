package public

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

type apiControllerI18N struct{}

func newApiControllerI18N() *apiControllerI18N {
	return &apiControllerI18N{}
}

func (c *apiControllerI18N) Register(r chi.Router) {
	r.Route("/i18n", func(r chi.Router) {
		r.Get("/t", c.getT)
	})
}

func (c *apiControllerI18N) getT(w http.ResponseWriter, r *http.Request) {

	urlQuery := r.URL.Query()
	prefix := urlQuery.Get("prefix")
	format := urlQuery.Get("fmt")

	if format == "" {
		format = "json"
	}

	if format != "json" && format != "js" {
		apiWriteBadRequest(w, fmt.Errorf("unknown format value: %s", format))
		return
	}

	l := i18n.GetLocalizer(r.Context())
	lang := l.Lang()
	keys := l.TT(prefix)

	switch format {
	case "json":
		olhttp.NewAPIResponse(struct {
			Lang language.Tag      `json:"lang"`
			T    map[string]string `json:"t"`
		}{
			Lang: lang,
			T:    keys,
		}).Write(w)
	case "js":
		w.Header().Add("Content-Type", "application/javascript")
		w.Write([]byte("Object.assign(window.i18n = window.i18n || {},"))
		enc := json.NewEncoder(w)
		err := enc.Encode(keys)
		if err != nil {
			write500(w, r, err)
			return
		}
		w.Write([]byte(");"))
	default:
		panic("unknown format value: " + format)
	}

}
