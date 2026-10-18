package public

import (
	"net/http/httptest"
	"testing"
)

func TestSearchPageSize(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  uint
	}{
		{"", 20}, {"invalid", 20}, {"0", 1}, {"-10", 1}, {"1", 1}, {"20", 20}, {"50", 50}, {"100", 100}, {"101", 100},
	} {
		t.Run(tc.value, func(t *testing.T) {
			got := getBooksSearchRequest(httptest.NewRequest("GET", "/search?ps="+tc.value, nil), nil)
			if got.PageSize != tc.want {
				t.Fatalf("got %d, want %d", got.PageSize, tc.want)
			}
		})
	}
}

func TestSearchIndexFilters(t *testing.T) {
	got := getBooksSearchRequest(httptest.NewRequest("GET", "/search?rating=PG-13&views.min=10&readers.max=100&reviews.min=-1&popularity.min=1.25&weightedRating.max=NaN&createdAt.min=2026-01-01&updatedAt.max=2026-02-30", nil), nil)
	if got.Rating != "PG-13" || got.Views.Min.Value != 10 || got.Readers.Max.Value != 100 || got.Popularity.Min.Value != 1.25 || got.CreatedAt.Min != "2026-01-01" {
		t.Fatalf("missing index filters: %+v", got.BookSearchFilters)
	}
	if got.Reviews.Min.Valid || got.WeightedRating.Max.Valid || got.UpdatedAt.Max != "" {
		t.Fatal("invalid bounds were accepted")
	}
}
