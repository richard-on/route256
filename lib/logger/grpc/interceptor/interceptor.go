package interceptor

import (
	"context"
	"time"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor logs gRPC requests using Logger.
func UnaryServerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		now := time.Now()
		resp, err := handler(ctx, req)
		if err != nil {
			switch {
			case status.Code(err) == codes.Internal,
				status.Code(err) == codes.DataLoss,
				status.Code(err) == codes.Unavailable:
				log.ErrorGRPC(req, resp, info, now, err)
			case status.Code(err) == codes.ResourceExhausted,
				status.Code(err) == codes.Unknown:
				log.WarnGRPC(req, resp, info, now, err)
			default:
				log.DebugGRPC(req, resp, info, now, err)
			}

			return resp, err
		}

		log.DebugGRPC(req, resp, info, now)
		return resp, err
	}
}
