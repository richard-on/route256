package domain

import (
	"context"
)

// CreateOrder creates a new order for a user reserving ordered products in a warehouse.
func (d *Domain) CreateOrder(ctx context.Context, user int64) (int64, error) {

	items, err := d.CheckoutRepo.GetCartItems(ctx, user)
	if err != nil {
		return 0, err
	}

	orderInfo, err := d.OrderCreator.CreateOrder(ctx, user, items)
	if err != nil {
		return 0, err
	}

	err = d.CheckoutRepo.ClearCart(ctx, user)
	if err != nil {
		return 0, err
	}

	return orderInfo, nil
}
