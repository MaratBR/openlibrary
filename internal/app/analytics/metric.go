package analytics

import (
	"time"

	"go.uber.org/fx"
)

const (
	MetricViews        = "view"
	MetricSearchClicks = "searchclick"
	MetricImpressions  = "impression"
)

type MetricType string

type MetricRecord struct {
	Type       MetricType
	Value      float64
	BookID     int64
	Samples    int64
	OccurredAt time.Time
}

func NewMetricRecord(
	type_ MetricType,
	value float64,
	bookID int64,
	occurredAt time.Time,
) MetricRecord {
	return MetricRecord{
		Samples:    1,
		Value:      value,
		Type:       type_,
		BookID:     bookID,
		OccurredAt: occurredAt,
	}
}

var metricModule = fx.Module("ol_analytics",
	fx.Provide(
		newMetricRepository,
		fx.Private,
	),

	fx.Provide(
		newMetricService,
	),
)

// go2tsdef:generate
type MetricValue struct {
	Samples  int64   `json:"samples"`
	ValueSum float64 `json:"valueSum"`
}

func NewMetricValue(samples int64, valueSum float64) MetricValue {
	return MetricValue{
		Samples:  samples,
		ValueSum: valueSum,
	}
}

// go2tsdef:generate
type MetricValues struct {
	Total MetricValue `json:"total"`
	Year  MetricValue `json:"year"`
	Month MetricValue `json:"month"`
	Week  MetricValue `json:"week"`
	Day   MetricValue `json:"day"`
}
