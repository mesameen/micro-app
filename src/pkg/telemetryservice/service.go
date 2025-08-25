package telemetryservice

import (
	"context"
	"fmt"
	"os"

	"github.com/mesameen/micro-app/src/pkg/logger"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Repo is a telemetry provider
type Repo interface {
	GetServiceName() string
	Infof(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span)
	MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error)
	MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error)
}

type Service struct {
	serviceName string
	lp          *log.LoggerProvider
	mp          *metric.MeterProvider
	tp          *trace.TracerProvider
	log         *zap.SugaredLogger
	meter       otelmetric.Meter
	tracer      oteltrace.Tracer
}

func NewTelemetry(ctx context.Context, serviceName string, serviceVersion string) (*Service, error) {
	res := newResource(serviceName, serviceVersion)
	lp, err := newLoggerProvider(ctx, res)
	if err != nil {
		logger.Errorf("Failed to create logger provider. Error: %v", err)
		return nil, fmt.Errorf("failed to create logger provider: %w", err)
	}

	zapLogger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
			otelzap.NewCore(serviceName, otelzap.WithLoggerProvider(lp)),
		),
	)

	mp, err := newMeterProvider(ctx, res)
	if err != nil {
		logger.Errorf("Failed to create metrics provider. Error: %v", err)
		return nil, fmt.Errorf("failed to create metrics provider: %w", err)
	}
	tp, err := newTracerProvider(ctx, res)
	if err != nil {
		logger.Errorf("Failed to create tracer provider. Error: %v", err)
		return nil, fmt.Errorf("failed to create trace provider: %w", err)
	}
	return &Service{
		lp:          lp,
		mp:          mp,
		tp:          tp,
		log:         zapLogger.Sugar(),
		meter:       mp.Meter(serviceName),
		tracer:      tp.Tracer(serviceName),
		serviceName: serviceName,
	}, nil
}

// GetServiceName returns the name of the service
func (s *Service) GetServiceName() string {
	return s.serviceName
}

// Infof logs at INFO level
func (s *Service) Infof(ctx context.Context, format string, args ...interface{}) {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	s.log.With("trace_id", spanCtx.TraceID()).With("span_id", spanCtx.SpanID()).Infof(format, args)
}

// Errorf logs at ERROR level
func (s *Service) Errorf(ctx context.Context, format string, args ...interface{}) {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	s.log.With("trace_id", spanCtx.TraceID()).With("span_id", spanCtx.SpanID()).Errorf(format, args)
}

func (s *Service) Fatalf(ctx context.Context, format string, args ...interface{}) {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	s.log.With("trace_id", spanCtx.TraceID()).With("span_id", spanCtx.SpanID()).Fatalf(format, args...)
}

// TraceStart starts a new span with a given name. The span must be ended by calling End
func (s *Service) TraceStart(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return s.tracer.Start(ctx, name)
}

// MeterInt64Histogram creates a new int64 histogram metric.
func (s *Service) MeterInt64Histogram(metric Metric) (otelmetric.Int64Histogram, error) {
	histogram, err := s.meter.Int64Histogram(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create histogram. Error: %v", err)
	}
	return histogram, err
}

// MeterInt64UpDownCounter creates a new int64 up down counter metric
func (s *Service) MeterInt64UpDownCounter(metric Metric) (otelmetric.Int64UpDownCounter, error) {
	counter, err := s.meter.Int64UpDownCounter(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter. Error: %v", err)
	}
	return counter, err
}

// Shutdown shuts donw the logger, meter and tracer
func (s *Service) Shutdown(ctx context.Context) {
	s.lp.Shutdown(ctx)
	s.mp.Shutdown(ctx)
	s.tp.Shutdown(ctx)
}
