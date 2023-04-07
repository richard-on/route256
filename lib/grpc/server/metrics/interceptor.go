package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor handles metrics.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		RequestCounter.WithLabelValues(info.FullMethod).Inc()

		now := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(now)

		statusCode := status.Code(err).String()

		HistogramResponseTime.WithLabelValues(info.FullMethod, statusCode).Observe(elapsed.Seconds())
		ResponseCounter.WithLabelValues(info.FullMethod, statusCode).Inc()

		return resp, err
	}
}
