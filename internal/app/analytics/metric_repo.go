package analytics

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/store"
)

type MetricRepository interface {
	AddMetrics(ctx context.Context, tx store.DBTX, metric []MetricRecord) error
}
