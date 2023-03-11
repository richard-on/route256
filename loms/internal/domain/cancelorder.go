package domain

import "context"

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
