package templates

import "encoding/json"

func jsonString(s string) string {
	if s == "" {
		return "\"\""
	}
	b, _ := json.Marshal(s)
	return string(b)
}
