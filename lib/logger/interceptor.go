package logger

import (
	"context"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor logs gRPC requests using Logger.
//
// Note: this interceptor only logs errors for now. In the future, more methods must be added.
func UnaryServerInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		resp, err := handler(ctx, req)
		if err != nil {
			log.Errorf(err, "method: %v", info.FullMethod)
		}

		return resp, err
	}
}
