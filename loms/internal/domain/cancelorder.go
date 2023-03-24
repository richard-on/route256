package domain

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"gitlab.ozon.dev/rragusskiy/homework-1/lib/workerpool"
)

var (
	ErrOrderCancelled = errors.New("order does not exist or has already been cancelled")
	ErrStockNotExists = errors.New("warehouse or sku does not exist")
)

// CancelOrder cancels order, makes previously reserved products available.
func (d *Domain) CancelOrder(ctx context.Context, orderID int64) error {
	err := d.Transactor.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		err := d.LOMSRepo.CancelOrder(ctxTX, orderID)
		if err != nil {
			return err
		}

		skus, stocks, err := d.LOMSRepo.RemoveItemsFromReserved(ctxTX, orderID)
		if err != nil {
			return err
		}

		for i, sku := range skus {
			err = d.LOMSRepo.IncreaseStock(ctxTX, sku, stocks[i])
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// CancelUnpaidOrders cancels orders that are awaiting payment (status model.AwaitingPayment)
// for more than given paymentTimeout.
func (d *Domain) CancelUnpaidOrders(ctx context.Context, paymentTimeout time.Duration) []error {
	unpaidOrders, err := d.LOMSRepo.ListUnpaidOrders(ctx, paymentTimeout)
	if err != nil {
		return []error{err}
	}

	wp := workerpool.New[int64, struct{}](ctx, d.config.MaxPoolWorkers)

	wp.SubmitMany(func(ctx context.Context, id int64) (struct{}, error) {
		err = d.CancelOrder(ctx, id)
		if err != nil {
			// Error while cancelling one order must not affect cancelling all other orders,
			// so just write err in a channel and try to cancel other orders.
			return struct{}{}, fmt.Errorf("cancelling order %v: %w", id, err)
		}

		return struct{}{}, nil
	}, unpaidOrders)

	wp.Wait()

	var cancelErrors []error
	for _, res := range wp.GetResult() {
		if res.Err != nil {
			cancelErrors = append(cancelErrors, err)
		}
	}

	return cancelErrors
}

// MonitorUnpaid monitors unpaid orders at a given rate.
func (d *Domain) MonitorUnpaid(ctx context.Context, errChan chan error) {
	ticker := time.NewTicker(d.config.CancelInterval)
	// Start a separate goroutine to check and cancel unpaid orders.
	for {
		select {
		// Run cancelling on each tick.
		case <-ticker.C:
			errSlice := d.CancelUnpaidOrders(ctx, d.config.PaymentTimeout)
			if len(errSlice) > 0 {
				for _, cancelErr := range errSlice {
					errChan <- cancelErr
				}
			}
		case <-ctx.Done():
			close(errChan)
			return
		}
	}
}
