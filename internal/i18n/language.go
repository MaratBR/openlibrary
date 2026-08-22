package i18n

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

type lang struct {
	Lang        language.Tag
	Passthrough bool
}

func detectLanguage(
	r *http.Request,
	defaultLanguageTag language.Tag,
	queryParameter, cookie string,
) lang {
	queryParameterValue := r.URL.Query().Get(queryParameter)

	if queryParameterValue != "" {
		if strings.EqualFold(queryParameterValue, "IDENTITY") {
			return lang{Passthrough: true, Lang: defaultLanguageTag}
		}

		tag, err := language.Parse(queryParameterValue)
		if err == nil {
			return lang{Lang: tag}
		}
	}

	cookieValue, err := r.Cookie(cookie)
	if err == nil && cookieValue.Value != "" {
		tag, err := language.Parse(cookieValue.Value)
		if err == nil {
			return lang{Lang: tag}
		}
	}

	accept := r.Header.Get("Accept-Language")
	preferredLanguages, _, _ := language.ParseAcceptLanguage(accept)
	if len(preferredLanguages) > 0 {
		return lang{Lang: preferredLanguages[0]}
	} else {
		return lang{Lang: defaultLanguageTag}
	}
}

func setLanguage(w http.ResponseWriter, cookie string, lang language.Tag) {
	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; Path=/; Max-Age=%d;", cookie, lang.String(), 34560000))
}
