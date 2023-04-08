package loms

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CancelOrder cancels order, makes previously reserved products available.
func (l *LOMS) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error) {
	err := validateOrder(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = l.domain.CancelOrder(ctx, req.GetOrderId())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderCancelled):
			return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, domain.ErrStockNotExists):
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
