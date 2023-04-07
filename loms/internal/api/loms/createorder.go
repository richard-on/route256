package loms

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (l *LOMS) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	items := make([]model.Item, 0, len(req.Items))
	for i, item := range req.Items {
		err = validateSKU(item.GetSku())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%v at item with index %v", err.Error(), i)
		}

		err = validateCount(item.GetCount())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%v at item with index %v", err.Error(), i)
		}

		items = append(items, convert.ToModelItem(item))
	}

	id, err := l.domain.CreateOrder(ctx, req.GetUser(), items)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInsufficientStocks):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &loms.CreateOrderResponse{
		OrderId: id,
	}, nil
}
