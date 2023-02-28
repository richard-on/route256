package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// CreateOrder calls loms.CreateOrder on all items in user's cart using LOMS gRPC client.
func (c *Client) CreateOrder(ctx context.Context, user int64) (domain.OrderInfo, error) {
	// Example items
	items := []domain.Item{
		{
			SKU:   123456,
			Count: 1,
		},
		{
			SKU:   654321,
			Count: 2,
		},
	}

	itemsReq := make([]*loms.Item, 0, 2)
	for _, item := range items {
		itemsReq = append(itemsReq, &loms.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}
	orderReq := &loms.CreateOrderRequest{
		User:  user,
		Items: itemsReq,
	}

	resp, err := c.lomsClient.CreateOrder(ctx, orderReq)
	if err != nil {
		return domain.OrderInfo{}, err
	}

	return domain.OrderInfo{
		OrderID: resp.GetOrderId(),
	}, nil
}
