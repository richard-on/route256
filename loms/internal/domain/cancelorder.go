package domain

import "context"

// CancelOrder cancels order, makes previously reserved products available.
func (d *Domain) CancelOrder(ctx context.Context, orderID int64) error {
	// Blank business logic
	return nil
}
