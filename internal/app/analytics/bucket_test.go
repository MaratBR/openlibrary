package analytics

import (
	"testing"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

func TestBucketStartsAllUsesUTCBoundaries(t *testing.T) {
	input := time.Date(2026, time.March, 8, 23, 45, 0, 0, time.FixedZone("UTC-7", -7*60*60))
	total, year, month, week, day := newBucketStartTime(input).All()

	assertTimeEqual(t, "total", total, time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC))
	assertTimeEqual(t, "year", year, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC))
	assertTimeEqual(t, "month", month, time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC))
	assertTimeEqual(t, "week", week, time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC))
	assertTimeEqual(t, "day", day, time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC))
}

func TestBucketStartsWeekBeginsOnMonday(t *testing.T) {
	tests := []struct {
		name string
		day  time.Time
		want time.Time
	}{
		{
			name: "monday",
			day:  time.Date(2026, time.March, 9, 12, 0, 0, 0, time.UTC),
			want: time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "sunday",
			day:  time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC),
			want: time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, got, _ := newBucketStartTime(tt.day).All()
			assertTimeEqual(t, "week", got, tt.want)
		})
	}
}

func TestBucketsWithLookbackIncludesCrossedPeriods(t *testing.T) {
	now := time.Date(2026, time.January, 1, 1, 0, 0, 0, time.UTC)
	buckets := newBucketStartTime(now).BucketsWithLookback(2 * time.Hour)

	wants := []bucketID{
		{Type: store.OlAnalyticsBucketPeriodTypeAll, Start: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{Type: store.OlAnalyticsBucketPeriodTypeYear, Start: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{Type: store.OlAnalyticsBucketPeriodTypeMonth, Start: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{Type: store.OlAnalyticsBucketPeriodTypeWeek, Start: time.Date(2025, time.December, 29, 0, 0, 0, 0, time.UTC)},
		{Type: store.OlAnalyticsBucketPeriodTypeDay, Start: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{Type: store.OlAnalyticsBucketPeriodTypeYear, Start: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{Type: store.OlAnalyticsBucketPeriodTypeMonth, Start: time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)},
		{Type: store.OlAnalyticsBucketPeriodTypeDay, Start: time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC)},
	}

	if len(buckets) != len(wants) {
		t.Fatalf("BucketsWithLookback() returned %d buckets, want %d: %#v", len(buckets), len(wants), buckets)
	}
	for i := range wants {
		if buckets[i].Type != wants[i].Type || !buckets[i].Start.Equal(wants[i].Start) {
			t.Errorf("bucket[%d] = %#v, want %#v", i, buckets[i], wants[i])
		}
	}
}

func TestBucketsWithLookbackWithinSameDayAddsNothing(t *testing.T) {
	now := time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC)
	if got, want := len(newBucketStartTime(now).BucketsWithLookback(2*time.Hour)), 5; got != want {
		t.Fatalf("BucketsWithLookback() returned %d buckets, want %d", got, want)
	}
}

func assertTimeEqual(t *testing.T, name string, got, want time.Time) {
	t.Helper()
	if !got.Equal(want) {
		t.Errorf("%s = %s, want %s", name, got, want)
	}
}
