package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProviderConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "provider_connections_total",
		Help: "Total number of connections to each provider",
	}, []string{"provider"})

	ProviderErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "provider_errors_total",
		Help: "Total number of provider errors",
	}, []string{"provider", "error_type"})

	TokenExchanges = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "token_exchanges_total",
		Help: "Total number of token exchanges",
	}, []string{"provider", "status"})

	UserSessions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "user_sessions_total",
		Help: "Total number of user sessions",
	}, []string{"provider"})

	AuthFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_failures_total",
		Help: "Total number of authentication failures",
	}, []string{"provider", "failure_reason"})

	ProviderLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "provider_request_duration_seconds",
		Help:    "Time spent making requests to external providers",
		Buckets: prometheus.DefBuckets,
	}, []string{"provider", "operation"})
)

func RecordProviderConnection(provider string) {
	ProviderConnections.WithLabelValues(provider).Inc()
}

func RecordProviderError(provider string, errorType string) {
	ProviderErrors.WithLabelValues(provider, errorType).Inc()
}

func RecordTokenExchange(provider string, status string) {
	TokenExchanges.WithLabelValues(provider, status).Inc()
}

func RecordUserSession(provider string) {
	UserSessions.WithLabelValues(provider).Inc()
}

func RecordAuthFailure(provider string, reason string) {
	AuthFailures.WithLabelValues(provider, reason).Inc()
}

func RecordProviderLatency(provider string, operation string, duration float64) {
	ProviderLatency.WithLabelValues(provider, operation).Observe(duration)
}
