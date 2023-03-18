package domain

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/workerpool"
)

// ListCart lists all products that are currently in a user's cart.
func (d *Domain) ListCart(ctx context.Context, user int64) ([]model.Item, uint32, error) {
	items, err := d.CheckoutRepo.GetCartItems(ctx, user)
	if err != nil {
		return nil, 0, err
	}

	// Create a child context that may be cancelled in a worker pool.
	poolCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wp := workerpool.New[model.Item, model.Item](poolCtx, d.config.MaxPoolWorkers)

	wp.SubmitMany(func(poolCtx context.Context, item model.Item) (model.Item, error) {
		product, err := d.ProductLister.GetProduct(poolCtx, item.SKU)
		if err != nil {
			// Even a single error will render the result unusable,
			// so cancel context and stop worker pool as soon as possible.
			cancel()
			return model.Item{}, err
		}

		item.ProductInfo = model.ProductInfo{
			Name:  product.Name,
			Price: product.Price,
		}

		return item, err

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
