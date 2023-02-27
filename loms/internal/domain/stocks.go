package domain

import "context"

type Stock struct {
	WarehouseID int64
	Count       uint64
}

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
