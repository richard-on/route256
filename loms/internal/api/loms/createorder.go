package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (l *LOMS) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {

	itemDomain := make([]domain.Item, 0, len(req.Items))
	for _, item := range req.Items {
		itemDomain = append(itemDomain, domain.Item{
			SKU:   item.GetSku(),
			Count: uint16(item.GetCount()),
		})
	}

	id, err := l.domain.CreateOrder(ctx, req.GetUser(), itemDomain)
	if err != nil {
		return nil, err
	}

	return &loms.CreateOrderResponse{
		OrderId: id,
	}, nil
}
