package templates

import (
	"time"

	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/a-h/templ"
	"github.com/eduardolat/goeasyi18n"
)

func i18nExtractKeys(l *i18n.Localizer, keys []string) templ.ComponentScript {
	m := make(map[string]string, len(keys))
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = l.T(keys[i])
	}

	return i18nKeys(m)
}

func _t(l *i18n.Localizer, key string) string {
	value := l.T(key)
	return value
}

func _tt(l *i18n.Localizer, key string, params any) string {
	value := l.T(key, goeasyi18n.Options{
		Data: params,
	})
	return value
}

func relativeTime(l *i18n.Localizer, t time.Time) string {
	dur := time.Now().Sub(t)
	s := dur.Seconds()
	if s < 60 {
		return l.T("time.justNow")
	}

	if s > 24*3600.0 {
		d := int(s / 3600.0 / 24.0)
		v := l.T("time.days", goeasyi18n.Options{
			Count: &d,
			Data:  map[string]interface{}{"count": d},
		})
		return v
	}

	if s > 3600.0 {
		h := int(s / 3600.0)
		v := l.T("time.hours", goeasyi18n.Options{
			Count: &h,
			Data:  map[string]interface{}{"count": h},
		})
		return v
	}

	m := int(s / 60.0)
	v := l.T("time.minutes", goeasyi18n.Options{
		Count: &m,
		Data:  map[string]interface{}{"count": m},
	})
	return v
}
