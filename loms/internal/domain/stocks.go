package domain

import "context"

// Stock represents a number of specific product available in a specific warehouse.
type Stock struct {
	// WarehouseID is the ID of a warehouse where the item is stored.
	WarehouseID int64
	// Count is the number of specific product available in this warehouse.
	Count uint64
}

// Stocks returns a number of available products with a given SKU in different warehouses.
func (d *Domain) Stocks(ctx context.Context, sku uint32) ([]Stock, error) {
	// Blank business logic
	return []Stock{
		{
			WarehouseID: 1,
			Count:       3,
		},
		{
			WarehouseID: 2,
			Count:       5,
		},
	}, nil
}
