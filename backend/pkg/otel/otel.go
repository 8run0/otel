package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
	ToolsCtxKey            = ContextKey("otel-tools")
)

type Tools struct {
	Tracer  trace.Tracer
	Meter   metric.Meter
	Cleanup func()
}

func GetToolsFromContext(ctx context.Context) *Tools {
	tools := ctx.Value(ToolsCtxKey).(*Tools)
	return tools
}

func NewTools(ctx context.Context, meterName string) (context.Context, *Tools) {
	// Registers a tracer Provider globally.
	cleanup := InstallExportPipeline(ctx, Resource())
	tracer := otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
	meter := global.MeterProvider().Meter(meterName)
	return ctx, &Tools{
		Tracer:  tracer,
		Meter:   meter,
		Cleanup: cleanup,
	}
}

func Resource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("otel-example"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
}

func InstallExportPipeline(ctx context.Context, resource *resource.Resource) func() {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("creating stdout exporter: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatalf("stopping tracer provider: %v", err)
		}
	}
}
