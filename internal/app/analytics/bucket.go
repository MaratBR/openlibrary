package analytics

import (
	"encoding/json"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type BucketStarts time.Time

func newBucketStartTime(day time.Time) BucketStarts {
	return BucketStarts(day)
}

func (bs BucketStarts) MarshalJSON() ([]byte, error) {
	total, year, month, week, day := bs.All()
	return json.Marshal(map[string]time.Time{
		"total": total,
		"year":  year,
		"month": month,
		"week":  week,
		"day":   day,
	})
}

func (bs BucketStarts) All() (time.Time, time.Time, time.Time, time.Time, time.Time) {
	year, month, day := time.Time(bs).UTC().Date()
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

func (bs BucketStarts) Buckets() []bucketID {

	total, year, month, week, day := bs.All()

	return []bucketID{
		{Type: store.OlAnalyticsBucketPeriodTypeAll, Start: total},
		{Type: store.OlAnalyticsBucketPeriodTypeYear, Start: year},
		{Type: store.OlAnalyticsBucketPeriodTypeMonth, Start: month},
		{Type: store.OlAnalyticsBucketPeriodTypeWeek, Start: week},
		{Type: store.OlAnalyticsBucketPeriodTypeDay, Start: day},
	}
}

func (bs BucketStarts) BucketsWithLookback(lookBackLimit time.Duration) []bucketID {
	ids := bs.Buckets()

	if lookBackLimit <= time.Duration(0) {
		return ids
	}

	past := BucketStarts(time.Time(bs).Add(-lookBackLimit))
	_, year, month, week, day := bs.All()
	_, year2, month2, week2, day2 := past.All()
	if year != year2 {
		ids = append(ids, bucketID{Type: store.OlAnalyticsBucketPeriodTypeYear, Start: year2})
	}
	if month != month2 {
		ids = append(ids, bucketID{Type: store.OlAnalyticsBucketPeriodTypeMonth, Start: month2})
	}
	if week != week2 {
		ids = append(ids, bucketID{Type: store.OlAnalyticsBucketPeriodTypeWeek, Start: week2})
	}
	if day != day2 {
		ids = append(ids, bucketID{Type: store.OlAnalyticsBucketPeriodTypeDay, Start: day2})
	}

	return ids
}

func (bs BucketStarts) Day() time.Time {
	_, _, _, _, day := bs.All()
	return day
}
