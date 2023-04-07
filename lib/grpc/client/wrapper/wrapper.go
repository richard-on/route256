package wrapper

import (
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/grpc/client/metrics"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/grpc/client/tracing"
	"google.golang.org/grpc"
)

// NewClient creates a new gRPC client wrapped with basic traces and metrics.
func NewClient(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			metrics.UnaryClientInterceptor(),
			tracing.UnaryClientInterceptor(opentracing.GlobalTracer()),
		),
	)

	return grpc.Dial(target, opts...)
}
