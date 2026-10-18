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

func TestSearchCacheKeyIncludesIndexFilters(t *testing.T) {
	baseline := getSearchRequestCacheKey(&BookSearchQuery{})
	for _, filters := range []BookSearchFilters{
		{Rating: "PG"}, {Views: SearchNumberRange{Min: Value(10.0)}},
		{Readers: SearchNumberRange{Max: Value(20.0)}}, {Reviews: SearchNumberRange{Min: Value(1.0)}},
		{Popularity: SearchNumberRange{Min: Value(1.25)}}, {WeightedRating: SearchNumberRange{Max: Value(4.5)}},
		{CreatedAt: SearchDateRange{Min: "2026-01-01"}}, {UpdatedAt: SearchDateRange{Max: "2026-02-01"}},
	} {
		if getSearchRequestCacheKey(&BookSearchQuery{BookSearchFilters: filters}) == baseline {
			t.Fatalf("cache ignores %+v", filters)
		}
	}
}

func TestSearchCacheKeyIncludesQuery(t *testing.T) {
	if getSearchRequestCacheKey(&BookSearchQuery{Query: "one"}) == getSearchRequestCacheKey(&BookSearchQuery{Query: "two"}) {
		t.Fatal("cache ignores text query")
	}
}

func TestActiveIndexFilters(t *testing.T) {
	query := DetailedBookSearchQuery{BookSearchFilters: BookSearchFilters{
		Rating: "G", Views: SearchNumberRange{Min: Value(0.0)}, Readers: SearchNumberRange{Max: Value(10.0)},
		CreatedAt: SearchDateRange{Min: "2026-01-01", Max: "2026-02-01"},
	}}
	if got := query.ActiveFilters(); got != 4 {
		t.Fatalf("got %d active filters, want 4", got)
	}
}
