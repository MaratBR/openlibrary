package app

import (
	"html"
	"strings"
)

// searchNameParts accepts only balanced highlight markers whose decoded text
// exactly reproduces the source title. Templates escape every text segment.
func searchNameParts(name, highlighted string) []SearchTextPart {
	if highlighted == "" {
		return nil
	}
	var parts []SearchTextPart
	var original strings.Builder
	matched := false
	for highlighted != "" {
		i := strings.IndexByte(highlighted, '<')
		if i < 0 {
			i = len(highlighted)
		}
		if i > 0 {
			text := html.UnescapeString(highlighted[:i])
			parts = append(parts, SearchTextPart{Text: text, Matched: matched})
			original.WriteString(text)
			highlighted = highlighted[i:]
		}
		if highlighted == "" {
			break
		}
		if !matched && strings.HasPrefix(highlighted, "<mark>") {
			matched = true
			highlighted = strings.TrimPrefix(highlighted, "<mark>")
		} else if matched && strings.HasPrefix(highlighted, "</mark>") {
			matched = false
			highlighted = strings.TrimPrefix(highlighted, "</mark>")
		} else {
			return nil
		}
	}
	if matched || original.String() != name {
		return nil
	}
	return parts
}
