package domain

import "context"

// Item represents a product to buy.
type Item struct {
	// SKU is the product's stock keeping unit.
	SKU uint32
	// Count is the number of product's with this SKU.
	Count uint16
}

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (d *Domain) CreateOrder(ctx context.Context, user int64, items []Item) (int64, error) {
	// Blank business logic
	return 42, nil
}
