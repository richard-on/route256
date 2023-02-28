package domain

import (
	"context"

	"github.com/pkg/errors"
)

var (
	// ErrInsufficientStocks is the error returned when current stock for an item is not enough to fulfill the order.
	ErrInsufficientStocks = errors.New("insufficient stock")
)

// Stock represents a number of specific product available in a specific warehouse.
type Stock struct {
	// WarehouseID is the ID of a warehouse where the item is stored.
	WarehouseID int64
	// Count is the number of specific product available in this warehouse.
	Count uint64
}

// AddToCart adds a number of items with given sku to user's cart.
func (d *Domain) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := d.stockChecker.Stocks(ctx, sku)
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
