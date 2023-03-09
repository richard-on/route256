package domain

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
)

var (
	// ErrNotInCart is the error returned when provided item is not found in user's cart.
	ErrNotInCart = errors.New("this item is not in cart")
	// ErrNotEnoughInCart is the error returned when cart count for this item is less than was requested to delete.
	ErrNotEnoughInCart = errors.New("cart count for this item is less than in delete request")
)

// DeleteFromCart deletes a number of items with given sku from user's cart.
func (d *Domain) DeleteFromCart(ctx context.Context, user int64, item model.Item) error {

	// Start transaction that reads cart count for this item and then decreases it or removes item for the database.
	err := d.Transactor.RunReadCommitted(ctx, func(ctxTX context.Context) error {
		count, err := d.CheckoutRepo.GetItemCartCount(ctxTX, user, item)
		if err != nil {
			return err
		}

		if count == int32(item.Count) {
			// If no items should be left in the cart after delete operation, delete the record entirely.
			if err = d.CheckoutRepo.DeleteItemCart(ctxTX, user, item.SKU); err != nil {
				return err
			}

		} else if count > int32(item.Count) {
			// Otherwise, decrease count for this item
			if err = d.CheckoutRepo.DecreaseItemCartCount(ctxTX, user, item); err != nil {
				return err
			}

		} else if count < int32(item.Count) {
			return errors.WithMessagef(ErrNotEnoughInCart, "item %v", item.SKU)
		} else {
			return ErrNotInCart
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
