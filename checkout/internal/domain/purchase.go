package domain

import (
	"context"
)

// OrderInfo represents information about newly created order.
type OrderInfo struct {
	// OrderID is the unique identifier ot this order.
	OrderID int64
}

// CreateOrder creates a new order for a user reserving ordered products in a warehouse.
func (d *Domain) CreateOrder(ctx context.Context, user int64) (OrderInfo, error) {
	orderInfo, err := d.orderCreator.CreateOrder(ctx, user)
	if err != nil {
		return OrderInfo{}, err
	}

	return orderInfo, nil
}
