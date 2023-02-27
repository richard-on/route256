package domain

import "context"

func (d *Domain) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	// Blank business logic
	return nil
}
