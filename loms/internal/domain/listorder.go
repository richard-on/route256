package domain

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

// ListOrder lists OrderInfo for a given orderID.
func (d *Domain) ListOrder(ctx context.Context, orderID int64) (model.Order, error) {

	var orderInfo model.Order
	err := d.Transactor.RunRepeatableRead(ctx, func(ctxTX context.Context) (err error) {

		orderInfo, err = d.LOMSRepo.ListOrderInfo(ctxTX, orderID)
		if err != nil {
			return err
		}

		orderInfo.Items, err = d.LOMSRepo.ListOrderItems(ctxTX, orderID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.Order{}, err
	}

	return orderInfo, nil
}
