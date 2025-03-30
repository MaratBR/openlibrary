package i18n

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/eduardolat/goeasyi18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

func tryParseAsObject(m map[string]any) (goeasyi18n.TranslateString, bool) {
	var (
		one  string
		zero string
		two  string
		few  string
		many string
		def  string
	)

	oneAny, ok := m["one"]
	if !ok {
		return goeasyi18n.TranslateString{}, false
	}
	one, ok = oneAny.(string)
	if !ok {
		return goeasyi18n.TranslateString{}, false
	}

	if many != "" {
		def = many
	} else {
		def = one
	}

	zeroAny, ok := m["zero"]
	if ok {
		zero, _ = zeroAny.(string)
	}

	twoAny, ok := m["two"]
	if ok {
		two, _ = twoAny.(string)
	}

	fewAny, ok := m["few"]
	if ok {
		few, _ = fewAny.(string)
	}

	manyAny, ok := m["many"]
	if ok {
		many, _ = manyAny.(string)
	}

	return goeasyi18n.TranslateString{
		Default: def,
		One:     one,
		Zero:    zero,
		Two:     two,
		Few:     few,
		Many:    many,
	}, true
}

func walkTranslations(
	fn func(goeasyi18n.TranslateString),
	prefix string,
	m map[string]any,
) {
	for k, v := range m {
		if vStr, ok := v.(string); ok {
			fn(goeasyi18n.TranslateString{
				Key:     prefix + k,
				Default: vStr,
			})
		} else if vMap, ok := v.(map[string]any); ok {
			ts, ok := tryParseAsObject(vMap)
			if ok {
				ts.Key = prefix + k
				fn(ts)
			} else {
				walkTranslations(fn, prefix+k+".", vMap)
			}
		}
	}
}

func loadFromTOML(files ...string) (goeasyi18n.TranslateStrings, error) {
	arr := make([]goeasyi18n.TranslateString, 0)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		m := make(map[string]any)

		err = toml.NewDecoder(f).Decode(&m)

		if err != nil {
			return nil, err
		}

		walkTranslations(func(ts goeasyi18n.TranslateString) {
			arr = append(arr, ts)
		}, "", m)
	}

	return goeasyi18n.TranslateStrings(arr), nil
}

func newI18N(
	lang language.Tag,
	files map[language.Tag][]string,
) *goeasyi18n.I18n {
	i18nInstance := goeasyi18n.NewI18n(goeasyi18n.Config{
		FallbackLanguageName:    lang.String(),
		DisableConsistencyCheck: false,
	})

	for lang, langFiles := range files {
		translations, err := loadFromTOML(langFiles...)
		if err != nil {
			panic(err)
		}

		i18nInstance.AddLanguage(lang.String(), translations)
	}

	return i18nInstance
}

type LocaleProvider struct {
	defaultLanguage language.Tag
	i18n            *goeasyi18n.I18n
	mx              sync.Mutex
	autoReload      bool
	files           map[language.Tag][]string
	lastLoad        time.Time
}

func NewLocaleProvider(
	defaultLanguage language.Tag,
	autoReload bool,
	files map[language.Tag][]string,
) *LocaleProvider {
	lp := &LocaleProvider{
		files:           files,
		autoReload:      autoReload,
		defaultLanguage: defaultLanguage,
	}
	lp.load()
	return lp
}

func (p *LocaleProvider) getI18N() *goeasyi18n.I18n {
	if p.autoReload && p.statChanged() {
		p.load()
	}

	return p.i18n
}

func (p *LocaleProvider) load() {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.i18n = newI18N(p.defaultLanguage, p.files)
	p.lastLoad = time.Now()
}

func (p *LocaleProvider) statChanged() bool {
	for _, files := range p.files {
		for _, file := range files {
			fileInfo, err := os.Stat(file)
			if err != nil {
				slog.Error("failed to stat i18n file", "file", file, "err", err)
				continue
			}

			if fileInfo.ModTime().After(p.lastLoad) {
				return true
			}
		}
	}

	return false
}

func (p *LocaleProvider) CreateLocalizer(r *http.Request) *Localizer {
	accept := r.Header.Get("Accept-Language")
	preferredLanguages, _, _ := language.ParseAcceptLanguage(accept)
	var tag language.Tag
	if len(preferredLanguages) > 0 {
		tag = preferredLanguages[0]
	} else {
		tag = p.defaultLanguage
	}
	i18nInstance := p.getI18N()
	localizer := newLocalizer(i18nInstance, tag)
	return localizer
}

type Localizer struct {
	i18n *goeasyi18n.I18n
	lang language.Tag
}

func (l *Localizer) T(key string, params ...goeasyi18n.Options) string {
	v := l.i18n.T(l.lang.String(), key, params...)
	if v == "" {
		return key
	}
	return v
}

func (l *Localizer) TData(key string, data any) string {
	return l.T(key, goeasyi18n.Options{
		Data: data,
	})
}

func newLocalizer(i18n *goeasyi18n.I18n, lang language.Tag) *Localizer {
	return &Localizer{
		i18n: i18n,
		lang: lang,
	}
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

func GetLocalizer(c context.Context) *Localizer {
	return c.Value(localizerKey).(*Localizer)
}
