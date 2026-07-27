package app

import "testing"

func TestParseBookSearchSort(t *testing.T) {
	tests := []struct {
		value string
		want  BookSearchSort
	}{
		{"", BookSearchSortRelevance},
		{"chapters", BookSearchSortChapters},
		{"words", BookSearchSortWords},
		{"words-per-chapter", BookSearchSortWordsPerChapter},
		{"created-at", BookSearchSortCreatedAt},
		{"last-updated", BookSearchSortLastUpdated},
		{"reviews", BookSearchSortReviews},
		{"readers", BookSearchSortReaders},
		{"weighted-rating", BookSearchSortWeightedRating},
		{"random", BookSearchSortRandom},
		{"unknown", BookSearchSortRelevance},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			if got := ParseBookSearchSort(test.value); got != test.want {
				t.Fatalf("ParseBookSearchSort(%q) = %q, want %q", test.value, got, test.want)
			}
		})
	}
}

func TestSearchCacheKeyIncludesSort(t *testing.T) {
	relevance := getSearchRequestCacheKey(&BookSearchQuery{})
	chapters := getSearchRequestCacheKey(&BookSearchQuery{Sort: BookSearchSortChapters})
	if relevance == chapters {
		t.Fatal("different search sorts produced the same cache key")
	}
}
