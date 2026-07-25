package templates

import (
	"fmt"
	"strconv"
)

func formatInt32(v int32) string {
	return fmt.Sprintf("%d", v)
}

func formatNumber(v int64) string {
	// Avoid negating v directly, which would overflow for math.MinInt64.
	abs := uint64(v)
	if v < 0 {
		abs = uint64(-(v + 1)) + 1
	}

	type unit struct {
		value  uint64
		suffix string
	}

	units := [...]unit{
		{1_000_000_000_000, "T"},
		{1_000_000_000, "B"},
		{1_000_000, "M"},
		{1_000, "K"},
	}

	for _, u := range units {
		if abs >= u.value {
			n := float64(abs) / float64(u.value)
			if v < 0 {
				n = -n
			}

			// One decimal at most; trailing ".0" is omitted.
			return strconv.FormatFloat(n, 'f', 1, 64)[:trimDecimalZero(n)] + u.suffix
		}
	}

	return strconv.FormatInt(v, 10)
}

func trimDecimalZero(v float64) int {
	s := strconv.FormatFloat(v, 'f', 1, 64)
	if len(s) >= 2 && s[len(s)-2:] == ".0" {
		return len(s) - 2
	}
	return len(s)
}
