package domain

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
)

// DeleteFromCart deletes a number of items with given sku from user's cart.
func (d *Domain) DeleteFromCart(ctx context.Context, user int64, item model.Item) error {

	err := d.Transactor.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		count, err := d.CheckoutRepo.GetItemCartCount(ctxTX, user, item)
		if err != nil {
			return err
		}

		if count-int32(item.Count) == 0 {
			if err = d.CheckoutRepo.DeleteItemCart(ctxTX, user, item.SKU); err != nil {
				return err
			}

		} else if count-int32(item.Count) > 0 {
			if err = d.CheckoutRepo.DecreaseItemCartCount(ctxTX, user, item); err != nil {
				return err
			}

		} else {
			return errors.New("this item is not in cart")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
