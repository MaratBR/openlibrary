package commonutil

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printerEnglish = message.NewPrinter(language.English)

func FormatInt(v int) string {
	return printerEnglish.Sprintf("%d", v)
}
