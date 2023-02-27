package domain

import "context"

// DeleteFromCart deletes a number of items with given sku from user's cart.
func (d *Domain) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	// Blank business logic
	return nil
}
