package telemetryservice

// Metric represents a metric that can be collected by the server
type Metric struct {
	Name        string
	Unit        string
	Description string
}

// MetricRequestDurationMilliSec is a metric that measures the latency of HTTP requests processed by server, in mulli seconds
var MetricRequestDurationMilliSec = Metric{
	Name:        "request_duration_millisec",
	Unit:        "ms",
	Description: "Measures the latency of HTTP requests processed by the server, in milliseconds",
}

// MetricRequestInFlight is a metric that measures the no of requests currently being process by server
var MetricRequestInFlight = Metric{
	Name:        "requests_inflight",
	Unit:        "{count}",
	Description: "Measures the no of requests currently handling by server",
}

var CounterMetric = Metric{
	Name: "api_calls",
	Unit: "{count}",
}
