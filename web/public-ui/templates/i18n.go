package templates

import (
	"github.com/a-h/templ"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func i18nExtractKeys(l *i18n.Localizer, keys []string) templ.ComponentScript {
	m := make(map[string]string, len(keys))
	msg := &i18n.LocalizeConfig{}
	for i := 0; i < len(keys); i++ {
		msg.MessageID = keys[i]
		m[keys[i]] = l.MustLocalize(msg)
	}

	return i18nKeys(m)
}
