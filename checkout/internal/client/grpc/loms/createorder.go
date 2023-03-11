package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// CreateOrder calls loms.CreateOrder on all items in user's cart using LOMS gRPC client.
func (c *Client) CreateOrder(ctx context.Context, user int64, items []model.Item) (int64, error) {

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
		return 0, err
	}

	return resp.GetOrderId(), nil
}
