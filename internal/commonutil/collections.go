package commonutil

import (
	"cmp"
	"slices"
)

func MergeArrays[T cmp.Ordered](a, b []T) []T {
	c := []T{}
	seen := map[T]struct{}{}

	for _, v := range a {
		if _, ok := seen[v]; !ok {
			c = append(c, v)
		}
	}

	for _, v := range b {
		if _, ok := seen[v]; !ok {
			c = append(c, v)
		}
	}

	slices.Sort(c)

	return c
}

func MapSlice[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
