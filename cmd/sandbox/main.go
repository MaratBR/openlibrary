package main

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
)

func main() {
	p := bluemonday.NewPolicy()
	p.AllowElements("olwidget", "ol-widget")
	p.AllowNoAttrs().OnElements("olwidget", "ol-widget")

	for _, s := range []string{
		`<olwidget name="tet">test</olwidget>`,
		`<ol-widget>test</ol-widget>`,
	} {
		fmt.Printf("%q => %q\n", s, p.Sanitize(s))
	}
}
