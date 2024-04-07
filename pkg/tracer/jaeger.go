package tracer

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func (c *Config) SetTracer(ctx context.Context) error {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(c.Host),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		return errors.Wrap(err, "create jaeger exporter")
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(c.ServiceName),
			),
		),
		trace.WithSampler(
			trace.TraceIDRatioBased(c.Ratio),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)

	return nil
}
