package htmlsanitizer

import "github.com/microcosm-cc/bluemonday"

var (
	htmlPolicy *bluemonday.Policy
)

func init() {
	htmlPolicy = bluemonday.UGCPolicy()
}

func Sanitize(s string) string {
	return htmlPolicy.Sanitize(s)
}
