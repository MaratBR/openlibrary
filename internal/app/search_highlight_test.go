package app

import (
	"reflect"
	"testing"
)

func TestSearchNameParts(t *testing.T) {
	for _, tc := range []struct {
		name, title, highlight string
		want                   []SearchTextPart
	}{
		{"multiple matches", "The Hobbit Returns", "The <mark>Hobbit</mark> <mark>Returns</mark>", []SearchTextPart{{Text: "The "}, {Text: "Hobbit", Matched: true}, {Text: " "}, {Text: "Returns", Matched: true}}},
		{"unicode and HTML", "魔法 <mark> & <script>", "<mark>魔法</mark> &lt;mark&gt; &amp; &lt;script&gt;", []SearchTextPart{{Text: "魔法", Matched: true}, {Text: " <mark> & <script>"}}},
		{"missing", "Book", "", nil},
		{"different source", "Book", "<mark>Other</mark>", nil},
		{"unclosed", "Book", "<mark>Book", nil},
		{"nested", "Book", "<mark><mark>Book</mark></mark>", nil},
		{"unexpected markup", "Book", "<script>Book</script>", nil},
		{"unexpected closing", "Book", "Book</mark>", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := searchNameParts(tc.title, tc.highlight); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}
