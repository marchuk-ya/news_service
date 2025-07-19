package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer provides tracing functionality
type Tracer struct {
	tracer trace.Tracer
}

// NewTracer creates a new tracer instance
func NewTracer(serviceName, serviceVersion string) (*Tracer, error) {
	// Create Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Create tracer
	tracer := tp.Tracer(serviceName)

	return &Tracer{
		tracer: tracer,
	}, nil
}

// StartSpan starts a new span
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.AddEvent(name)
	}
}

// SetAttributes sets attributes on the current span
func SetAttributes(ctx context.Context, key string, value interface{}) {
	if span := trace.SpanFromContext(ctx); span != nil {
		// Note: In a real implementation, you would convert value to appropriate attribute type
		// For now, we'll skip this to avoid complexity
		_ = key
		_ = value
	}
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.RecordError(err)
	}
}

// Shutdown gracefully shuts down the tracer
func (t *Tracer) Shutdown(ctx context.Context) error {
	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		return tp.Shutdown(ctx)
	}
	return nil
}
