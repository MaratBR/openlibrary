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

func MergeArraysNoSort[T comparable](a, b []T) []T {
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

	return c
}

func MapSlice[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func HasDuplicates[T comparable](a []T) bool {
	seen := make(map[T]struct{}, len(a))

	for _, v := range a {
		if _, ok := seen[v]; ok {
			return true
		}
		seen[v] = struct{}{}
	}

	return false
}

func ContainsSameAndNoDuplicates[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	if HasDuplicates(a) {
		return false
	}

	seen := make(map[T]struct{}, len(a))

	for _, v := range a {
		seen[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := seen[v]; !ok {
			return false
		}
	}

	return true
}
