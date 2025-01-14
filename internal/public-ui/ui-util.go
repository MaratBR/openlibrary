package publicui

import (
	"net/http"
	"strconv"
)

type uiSettings struct {
	FontSize int
}

func getUIBookSettings(r *http.Request) uiSettings {
	var settings uiSettings
	settings.FontSize = 18

	v, err := r.Cookie("ifs")
	if err == nil {
		intValue, err := strconv.Atoi(v.Value)
		if err == nil && (intValue >= 10 || intValue <= 99) {
			settings.FontSize = intValue
		}
	}

	return settings
}
