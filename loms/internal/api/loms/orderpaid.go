package loms

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// OrderPaid marks order as paid.
func (l *LOMS) OrderPaid(ctx context.Context, req *loms.OrderPaidRequest) (*emptypb.Empty, error) {
	err := validateOrder(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = l.domain.OrderPaid(ctx, req.GetOrderId())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotExistsOrPaid):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
