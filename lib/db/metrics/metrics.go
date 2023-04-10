package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "postgres",
		Name:      "requests_total",
	},
		[]string{"operation", "sql_query"},
	)

	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "postgres",
		Name:      "responses_total",
	},
		[]string{"operation", "sql_query", "sql_operation", "error"},
	)

	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "postgres",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"operation", "sql_query", "sql_operation", "error"},
	)

	TXCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "postgres",
		Name:      "tx_total",
	},
		[]string{"operation"},
	)
)
