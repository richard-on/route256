package domain

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
)

// ListCart lists all products that are currently in a user's cart.
func (d *Domain) ListCart(ctx context.Context, user int64) ([]model.Item, uint32, error) {

	items, err := d.CheckoutRepo.GetCartItems(ctx, user)
	if err != nil {
		return nil, 0, err
	}

	var totalPrice uint32 = 0
	for i, item := range items {
		product, err := d.ProductLister.GetProduct(ctx, item.SKU)
		if err != nil {
			return nil, 0, err
		}
		items[i].ProductInfo = model.ProductInfo{
			Name:  product.Name,
			Price: product.Price,
		}

		totalPrice += product.Price * uint32(item.Count)
	}

	return items, totalPrice, nil
}
