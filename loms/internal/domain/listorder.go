package domain

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

var (
	ErrEmptyOrder = errors.New("order does not exist")
)

// ListOrder lists OrderInfo for a given orderID.
func (d *Domain) ListOrder(ctx context.Context, orderID int64) (model.Order, error) {

	var orderInfo model.Order
	err := d.Transactor.RunReadCommitted(ctx, func(ctxTX context.Context) (err error) {

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
