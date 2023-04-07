package tracing

import (
	"context"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/tracing/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor handles gRPC server tracing.
func UnaryServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		spanContext, _ := utils.ExtractSpanContext(ctx, tracer)

		span := tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)
		span.LogFields(log.Object("request", req))

		resp, err := handler(ctx, req)
		if err != nil {
			otgrpc.SetSpanTags(span, err, false)
			span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		} else {
			span.SetTag("response_code", status.Code(err))
			span.LogFields(log.Object("response", resp))
		}

		return resp, err
	}
}
