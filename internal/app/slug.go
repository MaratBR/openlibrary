package app

import (
	"strings"

	"github.com/gosimple/slug"
)

func makeSlug(name string) string {
	slug := slug.Make(strings.ToLower(name))
	if len(slug) > 80 {
		slug = slug[:80]
	}
	return slug
}
