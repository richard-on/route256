package domain

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

var (
	ErrNotExistsOrPaid = errors.New("order does not exist or have already been paid")
)

// OrderPaid marks order as paid.
func (d *Domain) OrderPaid(ctx context.Context, orderID int64) error {

	// This transaction ensures that if order is paid, items are removed from reserve.
	err := d.Transactor.RunReadCommitted(ctx, func(ctxTX context.Context) (err error) {

		err = d.LOMSRepo.PayOrder(ctxTX, orderID)
		if err != nil {
			return err
		}

		_, _, err = d.LOMSRepo.RemoveItemsFromReserved(ctxTX, orderID)
		if err != nil {
			return err
		}

		err = d.CreateStatusMessage(ctx, orderID, model.Paid)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
