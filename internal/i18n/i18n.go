package i18n

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

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
	def *keyDef,
) {
	for k, v := range m {
		defInner := def.addInner(k)

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
				walkTranslations(fn, prefix+k+".", vMap, defInner)
			}
		}
	}
}

func loadFromTOML(files ...string) (goeasyi18n.TranslateStrings, *keyDef, error) {
	arr := make([]goeasyi18n.TranslateString, 0)

	fullDef := &keyDef{}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, nil, err
		}
		defer f.Close()

		m := make(map[string]any)

		err = toml.NewDecoder(f).Decode(&m)

		if err != nil {
			return nil, nil, err
		}

		def := &keyDef{}
		walkTranslations(func(ts goeasyi18n.TranslateString) {
			arr = append(arr, ts)
		}, "", m, def)
		fullDef.mergeWith(def)
	}

	return goeasyi18n.TranslateStrings(arr), fullDef, nil
}

func newI18N(
	lang language.Tag,
	files map[language.Tag][]string,
	log *zap.SugaredLogger,
) (*goeasyi18n.I18n, *keyDef) {
	i18nInstance := goeasyi18n.NewI18n(goeasyi18n.Config{
		FallbackLanguageName:    lang.String(),
		DisableConsistencyCheck: false,
	})

	allLangDef := &keyDef{}

	for lang, langFiles := range files {
		translations, langKeyDef, err := loadFromTOML(langFiles...)
		if err != nil {
			log.Errorw("failed to open language files", "files", langFiles, "lang", lang, "err", err)
			continue
		}
		allLangDef.mergeWith(langKeyDef)
		i18nInstance.AddLanguage(lang.String(), translations)
	}

	return i18nInstance, allLangDef
}

type LocaleProvider struct {
	defaultLanguage language.Tag
	i18n            *goeasyi18n.I18n
	def             *keyDef
	mx              sync.Mutex
	autoReload      bool
	files           map[language.Tag][]string
	lastLoad        time.Time
	queryParam      string
	cookie          string
	log             *zap.SugaredLogger
}

func NewLocaleProvider(
	defaultLanguage language.Tag,
	autoReload bool,
	files map[language.Tag][]string,
	log *zap.SugaredLogger,
) *LocaleProvider {
	lp := &LocaleProvider{
		files:           files,
		autoReload:      autoReload,
		defaultLanguage: defaultLanguage,
		queryParam:      "lang",
		cookie:          "lang",
		log:             log,
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
	p.i18n, p.def = newI18N(p.defaultLanguage, p.files, p.log)
	p.lastLoad = time.Now()
}

func (p *LocaleProvider) statChanged() bool {
	for _, files := range p.files {
		for _, file := range files {
			fileInfo, err := os.Stat(file)
			if err != nil {
				zap.S().Errorw("failed to stat i18n file", "file", file, "err", err)
				continue
			}

			if fileInfo.ModTime().After(p.lastLoad) {
				return true
			}
		}
	}

	return false
}

func (p *LocaleProvider) createLocalizer(r *http.Request) *Localizer {
	lang := detectLanguage(r, p.defaultLanguage, p.queryParam, p.cookie)
	i18nInstance := p.getI18N()
	localizer := newLocalizer(i18nInstance, p.def, lang.Lang, lang.Passthrough)
	return localizer
}

func (p *LocaleProvider) SetLanguage(w http.ResponseWriter, lang language.Tag) {
	setLanguage(w, p.cookie, lang)
}

type Localizer struct {
	i18n        *goeasyi18n.I18n
	def         *keyDef
	lang        language.Tag
	passthrough bool
}

func (l *Localizer) Lang() language.Tag {
	return l.lang
}

func (l *Localizer) Passthrough() bool {
	return l.passthrough
}

func (l *Localizer) T(key string, params ...goeasyi18n.Options) string {
	if l.passthrough {
		return key
	}

	v := l.i18n.T(l.lang.String(), key, params...)
	if v == "" {
		return key
	}
	return v
}

func (l *Localizer) TT(prefixPath string) map[string]string {
	o := l.def.path(prefixPath)
	if o == nil {
		return nil
	}

	m := make(map[string]string)

	if prefixPath != "" {
		prefixPath = prefixPath + "."
	}

	o.walkFullKeys(func(fullKey string) {
		m[prefixPath+fullKey] = l.T(prefixPath + fullKey)
	})

	return m
}

func (l *Localizer) TData(key string, data any) string {
	return l.T(key, goeasyi18n.Options{
		Data: data,
	})
}

func newLocalizer(i18n *goeasyi18n.I18n, def *keyDef, lang language.Tag, passthrough bool) *Localizer {
	return &Localizer{
		i18n:        i18n,
		lang:        lang,
		def:         def,
		passthrough: passthrough,
	}
}

type localizerKeyType struct{}

var localizerKey localizerKeyType

func (p *LocaleProvider) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := p.createLocalizer(r)
		r = r.WithContext(context.WithValue(r.Context(), localizerKey, l))
		next.ServeHTTP(w, r)
	})
}

func GetLocalizer(c context.Context) *Localizer {
	return c.Value(localizerKey).(*Localizer)
}

func TryGetLocalizer(c context.Context) (*Localizer, bool) {
	l := c.Value(localizerKey)
	if l == nil {
		return nil, false
	}
	return l.(*Localizer), true
}
