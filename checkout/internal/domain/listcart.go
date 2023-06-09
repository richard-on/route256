package domain

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/workerpool"
)

var (
	ErrEmptyCart = errors.New("cart is empty")
)

// ListCart lists all products that are currently in a user's cart.
func (d *Domain) ListCart(ctx context.Context, user int64) ([]model.Item, uint32, error) {
	items, err := d.CheckoutRepo.GetCartItems(ctx, user)
	if err != nil {
		return nil, 0, err
	}

	wp := workerpool.New[model.Item, model.Item](ctx, d.config.MaxPoolWorkers)

	wp.SubmitMany(func(ctx context.Context, item model.Item) (model.Item, error) {
		product, err := d.ProductLister.GetProduct(ctx, item.SKU)
		if err != nil {
			// Even a single error will render the result unusable,
			// so cancel context and stop worker pool as soon as possible.
			wp.StopNow()
			return model.Item{}, err
		}

		item.ProductInfo = model.ProductInfo{
			Name:  product.Name,
			Price: product.Price,
		}

		return item, nil

	}, items)

	wp.Wait()

	var totalPrice uint32 = 0
	for i, res := range wp.GetResult() {
		if res.Err != nil {
			return nil, 0, res.Err
		}
		items[i] = res.Value
		totalPrice += items[i].ProductInfo.Price * uint32(items[i].Count)
	}

	return items, totalPrice, nil
}
