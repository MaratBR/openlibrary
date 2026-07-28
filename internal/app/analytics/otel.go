package analytics

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/MaratBR/openlibrary/internal/app/analytics")
