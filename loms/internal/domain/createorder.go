package domain

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (d *Domain) CreateOrder(ctx context.Context, user int64, items []model.Item) (int64, error) {

	var orderID int64
	// This transaction inserts order info to database. After this step, the order is created.
	err := d.Transactor.RunReadCommitted(ctx, func(ctxTX context.Context) (err error) {
		orderID, err = d.LOMSRepo.InsertOrderInfo(ctxTX, model.Order{
			Status: model.NewOrder,
			User:   user,
		})
		if err != nil {
			return err
		}

		err = d.LOMSRepo.InsertOrderItems(ctxTX, orderID, items)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	// After order creation, this transaction checks whether we have enough items in stock to fulfill the order.
	// If it is fine, we pass ordered number of items into reserves.
	err = d.Transactor.RunRepeatableRead(ctx, func(ctxTX context.Context) (err error) {
		var stocks []model.Stock
		for _, item := range items {
			stocks, err = d.LOMSRepo.GetStocks(ctxTX, item.SKU)
			if err != nil {
				return err
			}

			toReserve := uint64(item.Count)
			for _, stock := range stocks {
				if toReserve < stock.Count {
					stock.Count = toReserve
				}
				toReserve -= stock.Count

				if err = d.LOMSRepo.DecreaseStock(ctxTX, int64(item.SKU), stock); err != nil {
					return errors.WithMessagef(err, "counting item with sku %v", item.SKU)
				}

				if err = d.LOMSRepo.ReserveItem(ctxTX, orderID, int64(item.SKU), stock); err != nil {
					return errors.WithMessagef(err, "reserving item with sku %v", item.SKU)
				}

				if toReserve == 0 {
					break
				}
			}

			if toReserve > 0 {
				return fmt.Errorf("order %v: sku %v: request %v items: not enough in stock",
					orderID, item.SKU, item.Count)
			}
		}

		// If we have successfully reserved items, change order status to "awaiting payment".
		err = d.LOMSRepo.ChangeOrderStatus(ctxTX, orderID, model.AwaitingPayment)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		// If we were not able to secure items, change order status to "failed".
		changeErr := d.LOMSRepo.ChangeOrderStatus(ctx, orderID, model.Failed)
		if changeErr != nil {
			return 0, changeErr
		}

		return 0, err
	}

	return orderID, nil
}
