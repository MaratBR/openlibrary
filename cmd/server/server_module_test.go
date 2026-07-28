package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestRouteTraceMiddlewareUsesRoutePattern(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previousProvider := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() { otel.SetTracerProvider(previousProvider) })

	inner := chi.NewRouter()
	inner.Get("/books/{bookID}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		router.Use(otelhttp.NewMiddleware(
			"HTTP request",
			otelhttp.WithTracerProvider(provider),
			otelhttp.WithSpanNameFormatter(routeSpanName),
		))
		router.Mount("/", inner)
	})

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/books/42", nil))

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
	if got, want := spans[0].Name(), "GET /books/{bookID}"; got != want {
		t.Fatalf("span name = %q, want %q", got, want)
	}
}
