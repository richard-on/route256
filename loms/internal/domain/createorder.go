package domain

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

// Item represents a product to buy.
type Item struct {
	// SKU is the product's stock keeping unit.
	SKU uint32
	// Count is the number of product's with this SKU.
	Count uint16
}

type Reserve struct {
	// WarehouseID is the ID of a warehouse where the item is reserved.
	WarehouseID int64
	// SKU is the product's stock keeping unit.
	SKU uint32
	// Count is the number of product's with this SKU.
	Count uint16
}

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (d *Domain) CreateOrder(ctx context.Context, user int64, items []Item) (int64, error) {

	var orderID int64
	err := d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) (err error) {
		orderID, err = d.Repository.InsertOrderInfo(ctxTX, OrderInfo{
			Status: NewOrder,
			User:   user,
		})
		if err != nil {
			return err
		}

		err = d.Repository.InsertOrderItems(ctxTX, orderID, items)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	err = d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) (err error) {
		var stocks []Stock
		for _, item := range items {
			stocks, err = d.Repository.GetStocks(ctxTX, item.SKU)
			if err != nil {
				return err
			}

			toReserve := uint64(item.Count)
			for _, stock := range stocks {
				if toReserve < stock.Count {
					stock.Count = toReserve
				}
				toReserve -= stock.Count

				if err = d.Repository.DecreaseStock(ctxTX, int64(item.SKU), stock); err != nil {
					return errors.WithMessagef(err, "counting item with sku %v", item.SKU)
				}

				if err = d.Repository.ReserveItem(ctxTX, orderID, int64(item.SKU), stock); err != nil {
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

		err = d.Repository.ChangeOrderStatus(ctx, orderID, AwaitingPayment)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		changeErr := d.Repository.ChangeOrderStatus(ctx, orderID, Failed)
		if changeErr != nil {
			return 0, changeErr
		}

		return 0, err
	}

	return orderID, nil
}
