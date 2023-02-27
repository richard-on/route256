package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// ListOrder lists order information.
func (l *LOMS) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {

	orderInfo, err := l.domain.ListOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	items := make([]*loms.Item, 0, len(orderInfo.Items))
	for _, item := range orderInfo.Items {
		items = append(items, &loms.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}

	return &loms.ListOrderResponse{
		Status: loms.Status(orderInfo.Status),
		User:   orderInfo.User,
		Items:  items,
	}, nil
}
