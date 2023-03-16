package domain

import (
	"context"
	"sync"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/workerpool"
)

const maxPoolWorkers int = 5

// ListCart lists all products that are currently in a user's cart.
func (d *Domain) ListCart(ctx context.Context, user int64) ([]model.Item, uint32, error) {

	items, err := d.CheckoutRepo.GetCartItems(ctx, user)
	if err != nil {
		return nil, 0, err
	}

	// Create a child context that may be cancelled in a worker pool.
	poolCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wp := workerpool.New(poolCtx, maxPoolWorkers)
	// errChan must be buffered as it is (theoretically) possible to get len(items) number of errors.
	errChan := make(chan error, len(items))
	mu := sync.Mutex{}

	var totalPrice uint32 = 0
	for i, item := range items {
		i := i
		item := item
		wp.Submit(func() {
			product, err := d.ProductLister.GetProduct(poolCtx, item.SKU)
			if err != nil {
				// Even a single error will render the result unusable,
				// so cancel context and stop worker pool as soon as possible.
				errChan <- err
				cancel()
				return
			}

			mu.Lock()
			items[i].ProductInfo = model.ProductInfo{
				Name:  product.Name,
				Price: product.Price,
			}

			totalPrice += product.Price * uint32(item.Count)
			mu.Unlock()
		})
	}
	wp.Wait()
	close(errChan)

	if err, ok := <-errChan; ok {
		return nil, 0, err
	}

	return items, totalPrice, nil
}
