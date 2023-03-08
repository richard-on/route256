package domain

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

// Stocks returns a number of available products with a given SKU in different warehouses.
func (d *Domain) Stocks(ctx context.Context, sku uint32) ([]model.Stock, error) {

	stocks, err := d.LOMSRepo.GetStocks(ctx, sku)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}
