package loms

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (l *LOMS) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, err
	}

	itemDomain := make([]domain.Item, 0, len(req.Items))
	for i, item := range req.Items {
		err = validateSKU(item.GetSku())
		if err != nil {
			return nil, errors.WithMessagef(err, "at item with index %d", i)
		}

		err = validateCount(item.GetCount())
		if err != nil {
			return nil, errors.WithMessagef(err, "at item with index %d", i)
		}

		singleItem := domain.Item{
			SKU:   item.GetSku(),
			Count: uint16(item.GetCount()),
		}
		itemDomain = append(itemDomain, singleItem)
	}

	id, err := l.domain.CreateOrder(ctx, req.GetUser(), itemDomain)
	if err != nil {
		return nil, err
	}

	return &loms.CreateOrderResponse{
		OrderId: id,
	}, nil
}
