package domain

import (
	"context"
)

const (
	Unspecified     = iota // Unspecified status.
	NewOrder               // NewOrder is the status for a newly created order.
	AwaitingPayment        // AwaitingPayment is the status for an order that awaits payment.
	Failed                 // Failed is the status for an order whose payment has failed.
	Paid                   // Paid is the status for a successfully paid order.
	Cancelled              // Cancelled is the status for a cancelled order.
)

// Status is an enumeration that represents a status of order payment.
type Status uint8

// OrderInfo represents information about the order, including its current Status,
// User who made the order and Items in the order.
type OrderInfo struct {
	// Status of order payment.
	Status Status
	// User
	User int64
	// Items is the slice of all Item's in user's cart
	Items []Item
}

// ListOrder lists OrderInfo for a given orderID.
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
