package domain

import (
	"context"

	"github.com/pkg/errors"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stock")
)

type Stock struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

func (d *Domain) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := d.stockChecker.Stock(ctx, sku)
	if err != nil {
		return errors.WithMessage(err, "checking stock")
	}

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			return nil
		}
	}

	return ErrInsufficientStocks
}
