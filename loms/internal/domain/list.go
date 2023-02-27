package domain

import (
	"context"
)

const (
	Unspecified = iota
	NewOrder
	AwaitingPayment
	Failed
	Payed
	Cancelled
)

type Status uint8

type OrderInfo struct {
	Status Status
	User   int64
	Items  []Item
}

func (d *Domain) ListOrder(ctx context.Context, orderID int64) (OrderInfo, error) {

	// Blank business logic

	return OrderInfo{
		Status: AwaitingPayment,
		User:   111111,
		Items: []Item{
			{
				SKU:   111111,
				Count: 5,
			},
			{
				SKU:   333333,
				Count: 12,
			},
		},
	}, nil
}
