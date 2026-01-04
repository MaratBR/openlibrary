package olhttp

import (
	"strconv"
	"strings"
)

func ParseInt64Slug(slug string) (int64, string) {
	// assume it's in slug form by default - "some-book-title-123"

	idx := strings.LastIndexByte(slug, '-')

	if idx == -1 {
		id, err := strconv.ParseInt(slug, 10, 64)
		if err == nil {
			return id, ""
		} else {
			return 0, ""
		}
	} else {
		lastPart := slug[idx+1:]
		id, err := strconv.ParseInt(lastPart, 10, 64)
		if err == nil {
			return id, slug[:idx]
		} else {
			return 0, ""
		}
	}
}
