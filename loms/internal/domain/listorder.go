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
	// Items is the slice of all Items in the order.
	Items []Item
}

// ListOrder lists OrderInfo for a given orderID.
func (d *Domain) ListOrder(ctx context.Context, orderID int64) (OrderInfo, error) {

	var orderInfo OrderInfo
	err := d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) (err error) {

		orderInfo, err = d.Repository.ListOrderInfo(ctxTX, orderID)
		if err != nil {
			return err
		}

		orderInfo.Items, err = d.Repository.ListOrderItems(ctxTX, orderID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return OrderInfo{}, err
	}

	return orderInfo, nil
}
