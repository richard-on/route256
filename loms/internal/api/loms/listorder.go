package loms

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListOrder lists order information.
func (l *LOMS) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {
	err := validateOrder(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orderInfo, err := l.domain.ListOrder(ctx, req.GetOrderId())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyOrder) || errors.Is(err, domain.ErrNoOrderItems):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &loms.ListOrderResponse{
		Status: loms.Status(orderInfo.Status),
		User:   orderInfo.User,
		Items:  convert.ToProtoItemSlice(orderInfo.Items),
	}, nil
}
