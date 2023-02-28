package domain

import "context"

// OrderPaid marks order as paid.
func (d *Domain) OrderPaid(ctx context.Context, orderID int64) error {
	// Blank business logic
	return nil
}
