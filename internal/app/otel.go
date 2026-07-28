package app

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/MaratBR/openlibrary/internal/app")

func startSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func endSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
