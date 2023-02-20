package cart

import (
	"context"

	"github.com/pkg/errors"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stock")
)

func (c *Check) Add(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := c.checker.Stock(ctx, sku)
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
