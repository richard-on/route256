package checkout

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
)

func (i *Implementation) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {

	// TODO: think about validation
	items, err := i.domain.ListCart(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	// TODO: total price should probably be calculated inside domain
	itemsResp := make([]*checkout.Item, 0, len(items))
	var totalPrice uint32 = 0
	for _, item := range items {
		itemsResp = append(itemsResp, &checkout.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
		totalPrice += item.Price
	}

	resp := checkout.ListCartResponse{
		Items:      itemsResp,
		TotalPrice: totalPrice,
	}

	return &resp, nil
}
