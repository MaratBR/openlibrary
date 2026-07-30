package analytics

import (
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type BucketStarts time.Time

func newBucketStartTime(day time.Time) BucketStarts {
	return BucketStarts(day)
}

func (bstu BucketStarts) All() (time.Time, time.Time, time.Time, time.Time, time.Time) {
	year, month, day := time.Time(bstu).UTC().Date()
	dayStart := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	weekDay := int(dayStart.Weekday())
	if weekDay == 0 {
		// sunday
		weekDay = 7
	}

	weekStart := dayStart.Add(time.Duration(24*-(weekDay-1)) * time.Hour)

	return time.Unix(0, 0).UTC(), time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(year, month, 1, 0, 0, 0, 0, time.UTC), weekStart, dayStart
}

type bucketID struct {
	Type  store.OlAnalyticsBucketPeriodType
	Start time.Time
}

func (bstu BucketStarts) Buckets() []bucketID {

	total, year, month, week, day := bstu.All()

	return []bucketID{
		{Type: store.OlAnalyticsBucketPeriodTypeAll, Start: total},
		{Type: store.OlAnalyticsBucketPeriodTypeYear, Start: year},
		{Type: store.OlAnalyticsBucketPeriodTypeMonth, Start: month},
		{Type: store.OlAnalyticsBucketPeriodTypeWeek, Start: week},
		{Type: store.OlAnalyticsBucketPeriodTypeDay, Start: day},
	}

}

func (bstu BucketStarts) Day() time.Time {
	_, _, _, _, day := bstu.All()
	return day
}
