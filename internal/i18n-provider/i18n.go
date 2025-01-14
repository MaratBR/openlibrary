package i18nprovider

import (
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

func loadI18NBundle(lang language.Tag) *i18n.Bundle {
	bundle := i18n.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFile("translations/en.toml")
	if err != nil {
		panic("failed to load en.toml: " + err.Error())
	}

	return bundle
}

type LocaleProvider struct {
}

func NewLocaleProvider() *LocaleProvider {
	return &LocaleProvider{}
}

func (p *LocaleProvider) GetLocalizer(r *http.Request) *i18n.Localizer {
	accept := r.Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(loadI18NBundle(language.English), accept)
	return localizer
}

type localizerKeyType struct{}

var localizerKey localizerKeyType

func (p *LocaleProvider) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := p.GetLocalizer(r)
		r = r.WithContext(context.WithValue(r.Context(), localizerKey, l))
		next.ServeHTTP(w, r)
	})
}

func GetLocalizer(c context.Context) *i18n.Localizer {
	return c.Value(localizerKey).(*i18n.Localizer)
}
