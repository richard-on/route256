package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor handles metrics.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		RequestCounter.WithLabelValues(method).Inc()

		now := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		elapsed := time.Since(now)

		statusCode := status.Code(err).String()

		HistogramResponseTime.WithLabelValues(method, statusCode).Observe(elapsed.Seconds())
		ResponseCounter.WithLabelValues(method, statusCode).Inc()

		return err
	}
}
