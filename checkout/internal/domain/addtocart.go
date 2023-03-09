package domain

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
)

var (
	// ErrInsufficientStocks is the error returned when current stock for an item is not enough to fulfill the order.
	ErrInsufficientStocks = errors.New("insufficient stock")
)

// AddToCart adds a number of items with given sku to user's cart.
func (d *Domain) AddToCart(ctx context.Context, user int64, item model.Item) error {
	stocks, err := d.StockChecker.Stocks(ctx, item.SKU)
	if err != nil {
		return errors.WithMessage(err, "checking stock")
	}

	count, err := d.CheckoutRepo.GetItemCartCount(ctx, user, item)
	if err != nil && !errors.Is(err, ErrNotInCart) {
		return err
	}

	inStock := false
	counter := int64(item.Count) + int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			inStock = true
			break
		}
	}

	if !inStock {
		return ErrInsufficientStocks
	}

	err = d.CheckoutRepo.AddToCart(ctx, user, item)
	if err != nil {
		return err
	}

	return nil
}
