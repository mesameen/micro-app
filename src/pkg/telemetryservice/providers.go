package telemetryservice

import (
	"context"
	"os"
	"time"

	"github.com/mesameen/micro-app/src/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewLoggerProvider initializes logger provider with the OTLP grpc exporter
func newLoggerProvider(ctx context.Context, res *resource.Resource) (*log.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint("localhost:4317"), otlploggrpc.WithInsecure())
	if err != nil {
		logger.Errorf("Failed to create otlploggrpc exporter. Error: %v", err)
		return nil, err
	}
	processor := log.NewBatchProcessor(exporter)
	lp := log.NewLoggerProvider(log.WithProcessor(processor), log.WithResource(res))
	return lp, nil
}

// NewMeterProvider creates a new meter provider with the OTLP grpc exporter
func newMeterProvider(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		logger.Errorf("Failed to create OTLP metric exporter: %v", err)
		return nil, err
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(5*time.Second))

	mp := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(res),
	)

	// Start Go runtime metrics collection
	if err := runtime.Start(
		runtime.WithMeterProvider(mp),
		runtime.WithMinimumReadMemStatsInterval(5*time.Second),
	); err != nil {
		logger.Errorf("Failed to start Go runtime instrumentation: %v", err)
		return nil, err
	}

	return mp, nil
}

// NewTracerProvider creates a new tracer provider with the OTLP grpc exporter
func newTracerProvider(ctx context.Context, res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
	if err != nil {
		logger.Errorf("Failed to create otlptracegrpc exporter. Error: %v", err)
		return nil, err
	}
	batcher := trace.WithBatcher(exporter)
	tp := trace.NewTracerProvider(batcher, trace.WithResource(res))
	return tp, nil
}

// NewResource creates a new OTEL resource withe the service name and version
func newResource(serviceName string, serviceVersion string) *resource.Resource {
	hostName, err := os.Hostname()
	if err != nil {
		logger.Errorf("failed to get hostname. Error: %v", err)
	}
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.HostName(hostName),
		attribute.String("service.name", serviceName),
	)
}
