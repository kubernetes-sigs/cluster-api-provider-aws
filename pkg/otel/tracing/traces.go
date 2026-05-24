package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

// function to initialize tracer
func InitTracer(ctx context.Context, traceSamplingRatio float64, tracingEndPoint string) (*sdktrace.TracerProvider, error) {

	if tracingEndPoint == "" {
		return nil, nil
	}
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(tracingEndPoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		slog.Error("failed to create OTLP trace exporter", "error", err)
		return nil, err
	}

	// Creating the resource for sdktrace
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("cluster-api-provider-aws"),
		semconv.ServiceVersionKey.String(version.Get().String()),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(
			sdktrace.ParentBased(
				sdktrace.TraceIDRatioBased(traceSamplingRatio),
			),
		),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

// function to start span for given context
func StartSpan(ctx context.Context, tracerName string, spanName string) (context.Context, trace.Span) {

	tr := otel.Tracer(tracerName)
	return tr.Start(ctx, spanName)
}
