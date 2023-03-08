package loms

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (l *LOMS) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, err
	}

	items := make([]model.Item, 0, len(req.Items))
	for i, item := range req.Items {

		err = validateSKU(item.GetSku())
		if err != nil {
			return nil, errors.WithMessagef(err, "at item with index %d", i)
		}

		err = validateCount(item.GetCount())
		if err != nil {
			return nil, errors.WithMessagef(err, "at item with index %d", i)
		}

		items = append(items, convert.ToModelItem(item))
	}

	id, err := l.domain.CreateOrder(ctx, req.GetUser(), items)
	if err != nil {
		return nil, err
	}

	return &loms.CreateOrderResponse{
		OrderId: id,
	}, nil
}
