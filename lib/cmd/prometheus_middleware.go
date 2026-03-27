package cmd

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "postal_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "postal_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "postal_http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	// Custom application metrics
	expandOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "postal_expand_operations_total",
			Help: "Total number of address expansion operations",
		},
		[]string{"status"},
	)

	parseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "postal_parse_operations_total",
			Help: "Total number of address parsing operations",
		},
		[]string{"status"},
	)

	expandAddressLength = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "postal_expand_address_length_bytes",
			Help:    "Length of addresses being expanded in bytes",
			Buckets: []float64{10, 50, 100, 200, 500, 1000},
		},
	)

	expandResultCount = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "postal_expand_result_count",
			Help:    "Number of expansion results returned",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100},
		},
	)

	parseAddressLength = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "postal_parse_address_length_bytes",
			Help:    "Length of addresses being parsed in bytes",
			Buckets: []float64{10, 50, 100, 200, 500, 1000},
		},
	)

	parseComponentCount = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "postal_parse_component_count",
			Help:    "Number of parsed components returned",
			Buckets: []float64{1, 2, 3, 5, 7, 10},
		},
	)
)

// PrometheusMiddleware returns a gin middleware for Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		start := time.Now()
		httpRequestsInFlight.Inc()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
		httpRequestsInFlight.Dec()
	}
}

// RecordExpandMetrics records metrics for expand operations
func RecordExpandMetrics(addressLen int, resultCount int, status string) {
	expandOperationsTotal.WithLabelValues(status).Inc()
	expandAddressLength.Observe(float64(addressLen))
	expandResultCount.Observe(float64(resultCount))
}

// RecordParseMetrics records metrics for parse operations
func RecordParseMetrics(addressLen int, componentCount int, status string) {
	parseOperationsTotal.WithLabelValues(status).Inc()
	parseAddressLength.Observe(float64(addressLen))
	parseComponentCount.Observe(float64(componentCount))
}
