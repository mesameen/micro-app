package telemetryservice

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
)

// LogRequest is a gin middleware that logs the request path
func (s *Service) LogRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.Infof(ctx.Request.Context(), "request to: %s", ctx.Request.URL.Path)
		ctx.Next()
		s.Infof(ctx.Request.Context(), "end of request to: %s", ctx.Request.URL.Path)
	}
}

func (s *Service) MeterRequestDuration() gin.HandlerFunc {
	histogram, err := s.MeterInt64Histogram(MetricRequestDurationMilliSec)
	if err != nil {
		s.Fatalf(context.Background(), "failed to create histogram: %w", err)
	}
	return func(ctx *gin.Context) {
		// capturing the request start time
		startTime := time.Now()
		ctx.Next()
		// record the request duration
		duration := time.Since(startTime)
		histogram.Record(
			ctx.Request.Context(),
			duration.Milliseconds(),
			metric.WithAttributes(
				httpconv.ServerRequest(s.GetServiceName(), ctx.Request)...,
			),
		)
	}
}

func (s *Service) MeterRequestsInFlight() gin.HandlerFunc {
	counter, err := s.MeterInt64UpDownCounter(MetricRequestInFlight)
	if err != nil {
		s.Fatalf(context.Background(), "failed to create counter: %w", err)
	}
	return func(ctx *gin.Context) {
		attrs := metric.WithAttributes(httpconv.ServerRequest(s.GetServiceName(), ctx.Request)...)
		// increases the number of requests in flight
		counter.Add(ctx.Request.Context(), 1, attrs)
		// execute next http handler
		ctx.Next()
		// decrease the number of requests in flight
		counter.Add(ctx.Request.Context(), -1, attrs)
	}
}
