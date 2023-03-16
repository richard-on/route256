package domain

import (
	"context"
	"fmt"
	"time"

	"gitlab.ozon.dev/rragusskiy/homework-1/lib/workerpool"
)

const maxPoolWorkers int = 5

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

	wp := workerpool.New(ctx, maxPoolWorkers)
	// errChan must be buffered as it is possible to get len(unpaidOrders) number of errors.
	errChan := make(chan error, len(unpaidOrders))

	for _, id := range unpaidOrders {
		// Cancel orders in a worker pool for efficiency.
		wp.Submit(func() {
			err = d.CancelOrder(ctx, id)
			if err != nil {
				// Error while cancelling one order must not affect cancelling all other orders,
				// so just write err in a channel and try to cancel other orders.
				errChan <- fmt.Errorf("cancelling order %v: %w", id, err)
			}
		})
	}
	// Wait for all orders to cancel.
	wp.Wait()
	close(errChan)

	// Return slice of all errors that may have occurred during cancelling.
	var cancelErrors []error
	for err = range errChan {
		cancelErrors = append(cancelErrors, err)
	}

	return cancelErrors
}
