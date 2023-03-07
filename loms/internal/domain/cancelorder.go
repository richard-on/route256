package domain

import "context"

// CancelOrder cancels order, makes previously reserved products available.
func (d *Domain) CancelOrder(ctx context.Context, orderID int64) error {

	err := d.Repository.CancelOrder(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}
