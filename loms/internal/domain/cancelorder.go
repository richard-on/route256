package domain

import "context"

// CancelOrder cancels order, makes previously reserved products available.
func (d *Domain) CancelOrder(ctx context.Context, orderID int64) error {

	err := d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		err := d.Repository.CancelOrder(ctxTX, orderID)
		if err != nil {
			return err
		}

		skus, stocks, err := d.Repository.RemoveItemsFromReserved(ctxTX, orderID)

		for i, sku := range skus {
			err = d.Repository.IncreaseStock(ctxTX, sku, stocks[i])
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
