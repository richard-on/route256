package domain

import (
	"context"
)

type OrderInfo struct {
	OrderID int64
}

func (d *Domain) CreateOrder(ctx context.Context, user int64) (OrderInfo, error) {
	orderInfo, err := d.orderCreator.Order(ctx, user)
	if err != nil {
		return OrderInfo{}, err
	}

	return orderInfo, nil
}
