package tech

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var currencyTr = otel.Tracer("currency")

func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return currencyTr.Start(ctx, spanName, opts...)
}
