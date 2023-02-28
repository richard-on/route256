package checkout

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
)

// ListCart lists all products that are currently in a user's cart.
func (c *Checkout) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, err
	}

	items, totalPrice, err := c.domain.ListCart(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	itemsResp := make([]*checkout.Item, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, &checkout.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
	}

	resp := checkout.ListCartResponse{
		Items:      itemsResp,
		TotalPrice: totalPrice,
	}

	return &resp, nil
}
