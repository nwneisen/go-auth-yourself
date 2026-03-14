package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	AuthEvents = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_events_total",
		Help: "Total number of authentication events",
	}, []string{"provider", "event_type", "status"})

	Sessions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "auth_sessions_active",
		Help: "Current number of active sessions",
	}, []string{"provider"})

	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"handler", "method", "status_code"})

	HTTPLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "method"})

	LastError = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "last_error_timestamp",
		Help: "Timestamp of the last error",
	})

	ProviderStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "provider_status",
		Help: "Health status of each provider (1=healthy, 0=unhealthy)",
	}, []string{"provider"})
)

func RecordAuthEvent(provider string, eventType string, status string) {
	AuthEvents.WithLabelValues(provider, eventType, status).Inc()
}

func RecordHTTPRequest(handler string, method string, statusCode int) {
	HTTPRequests.WithLabelValues(handler, method, statusCodeToString(statusCode)).Inc()
}

func statusCodeToString(code int) string {
	switch {
	case code >= 100 && code < 200:
		return "1xx"
	case code >= 200 && code < 300:
		return "2xx"
	case code >= 300 && code < 400:
		return "3xx"
	case code >= 400 && code < 500:
		return "4xx"
	case code >= 500 && code < 600:
		return "5xx"
	default:
		return "other"
	}
}

func RecordHTTPLatency(handler string, method string, duration float64) {
	HTTPLatency.WithLabelValues(handler, method).Observe(duration)
}

func SetLastError() {
	LastError.Set(float64(1))
}

func ClearLastError() {
	LastError.Set(float64(0))
}

func SetProviderHealth(provider string, healthy bool) {
	if healthy {
		ProviderStatus.WithLabelValues(provider).Set(1)
	} else {
		ProviderStatus.WithLabelValues(provider).Set(0)
	}
}
