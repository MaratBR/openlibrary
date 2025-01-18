package templates

import (
	"time"

	"github.com/a-h/templ"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func i18nExtractKeys(l *i18n.Localizer, keys []string) templ.ComponentScript {
	m := make(map[string]string, len(keys))
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = _t(l, keys[i])
	}

	return i18nKeys(m)
}

func _t(l *i18n.Localizer, key string) string {
	value, err := l.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		return key
	}
	return value
}

func _tt(l *i18n.Localizer, key string, params any) string {
	value, err := l.Localize(&i18n.LocalizeConfig{MessageID: key, TemplateData: params})
	if err != nil {
		return key
	}
	return value
}

func relativeTime(l *i18n.Localizer, t time.Time) string {
	dur := time.Now().Sub(t)
	s := dur.Seconds()
	if s < 60 {
		return _t(l, "time.justNow")
	}

	if s > 24*3600.0 {
		d := int(s / 3600.0 / 24.0)
		v, err := l.Localize(&i18n.LocalizeConfig{MessageID: "time.days", PluralCount: d, TemplateData: map[string]interface{}{"count": d}})
		if err != nil {
			return "ERROR:" + err.Error()
		}
		return v
	}

	if s > 3600.0 {
		h := int(s / 3600.0)
		v, err := l.Localize(&i18n.LocalizeConfig{MessageID: "time.hours", PluralCount: h, TemplateData: map[string]interface{}{"count": h}})
		if err != nil {
			return "ERROR:" + err.Error()
		}
		return v
	}

	m := int(s / 60.0)
	v, err := l.Localize(&i18n.LocalizeConfig{MessageID: "time.minutes", PluralCount: m, TemplateData: map[string]interface{}{"count": m}})
	if err != nil {
		return "ERROR:" + err.Error()
	}
	return v
}
