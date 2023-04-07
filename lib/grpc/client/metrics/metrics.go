package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc_client",
		Name:      "requests_total",
	},
		[]string{"method"},
	)

	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc_client",
		Name:      "responses_total",
	},
		[]string{"method", "status"},
	)

	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc_client",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"method", "status"},
	)
)
