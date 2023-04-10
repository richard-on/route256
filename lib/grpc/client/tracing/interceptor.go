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

// UnaryClientInterceptor handles gRPC client tracing.
func UnaryClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, resp interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}

		span := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		)
		defer span.Finish()

		ctx = utils.InjectSpanContext(ctx, tracer, span)

		err := invoker(ctx, method, req, resp, cc, opts...)
		if err != nil {
			otgrpc.SetSpanTags(span, err, true)
			span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		} else {
			span.SetTag("response_code", status.Code(err))
			span.LogFields(log.Object("response", resp))
		}

		return err
	}
}
