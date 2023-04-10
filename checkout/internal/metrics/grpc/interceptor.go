package grpc

import (
	"context"
	"time"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor handles metrics.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		metrics.RequestsCounter.WithLabelValues(info.FullMethod).Inc()

		now := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(now)

		statusCode := status.Code(err).String()

		metrics.HistogramResponseTime.WithLabelValues(info.FullMethod, statusCode).Observe(elapsed.Seconds())
		metrics.ResponseCounter.WithLabelValues(info.FullMethod, statusCode).Inc()

		return resp, err
	}
}
