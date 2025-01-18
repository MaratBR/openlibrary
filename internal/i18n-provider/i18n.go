package i18nProvider

import (
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

func LoadBundle(lang language.Tag, files ...string) *i18n.Bundle {
	bundle := i18n.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, file := range files {
		_, err := bundle.LoadMessageFile(file)
		if err != nil {
			panic("failed to load en.toml: " + err.Error())
		}
	}

	return bundle
}

type LocaleProvider struct {
	defaultLanguage language.Tag
	bundle          *i18n.Bundle
	mx              sync.Mutex
	autoReload      bool
	files           []string
	lastLoad        time.Time
}

func NewLocaleProvider(
	defaultLanguage language.Tag,
	autoReload bool,
	files []string,
) *LocaleProvider {
	lp := &LocaleProvider{
		files:      files,
		autoReload: autoReload,
	}
	lp.loadBundle()
	return lp
}

func (p *LocaleProvider) CreateLocalizer(r *http.Request) *i18n.Localizer {
	accept := r.Header.Get("Accept-Language")
	bundle := p.getBundle()
	localizer := i18n.NewLocalizer(bundle, accept)
	return localizer
}

func (p *LocaleProvider) getBundle() *i18n.Bundle {
	if p.autoReload && p.statChanges() {
		p.loadBundle()
	}

	return p.bundle
}

func (p *LocaleProvider) loadBundle() {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.bundle = LoadBundle(p.defaultLanguage, p.files...)
	p.lastLoad = time.Now()
}

func (p *LocaleProvider) statChanges() bool {
	for _, file := range p.files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			slog.Error("failed to stat i18n file", "file", file, "err", err)
			continue
		}

		if fileInfo.ModTime().After(p.lastLoad) {
			return true
		}
	}

	return false
}

type localizerKeyType struct{}

var localizerKey localizerKeyType

func (p *LocaleProvider) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := p.CreateLocalizer(r)
		r = r.WithContext(context.WithValue(r.Context(), localizerKey, l))
		next.ServeHTTP(w, r)
	})
}

func GetLocalizer(c context.Context) *i18n.Localizer {
	return c.Value(localizerKey).(*i18n.Localizer)
}
