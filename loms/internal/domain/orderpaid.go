package domain

import "context"

// OrderPaid marks order as paid.
func (d *Domain) OrderPaid(ctx context.Context, orderID int64) error {

	// This transaction ensures that if order is paid, items are removed from reserve.
	err := d.Transactor.RunReadCommitted(ctx, func(ctxTX context.Context) (err error) {

		err = d.LOMSRepo.PayOrder(ctx, orderID)
		if err != nil {
			return err
		}

		_, _, err = d.LOMSRepo.RemoveItemsFromReserved(ctx, orderID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
