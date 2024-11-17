package commonutil

import "slices"

func MergeStringArrays(a, b []string) []string {
	c := []string{}
	seen := map[string]struct{}{}

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
