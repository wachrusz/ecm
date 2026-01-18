package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics collects various Prometheus metrics.
type Metrics struct {
	RequestsTotal        *prometheus.CounterVec
	RequestLatency       *prometheus.HistogramVec
	CacheOperationsTotal prometheus.Counter
	CacheHitCounter      prometheus.Counter
	CacheMissCounter     prometheus.Counter
	PendingMessagesTotal prometheus.Gauge
}

// NewMetrics creates and initializes a Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "p2p_message_service_http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		RequestLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "p2p_message_service_http_request_duration_seconds",
				Help:    "HTTP request latency distribution",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
		CacheHitCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "p2p_message_service_cache_hits_total",
				Help: "Total number of cache hits",
			},
		),
		CacheMissCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "p2p_message_service_cache_misses_total",
				Help: "Total number of cache misses",
			},
		),
		CacheOperationsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "p2p_message_service_cache_operations_total",
				Help: "Total number of cache operations",
			},
		),
		PendingMessagesTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "p2p_message_service_pending_messages",
				Help: "Current number of pending messages",
			},
		),
	}
}
