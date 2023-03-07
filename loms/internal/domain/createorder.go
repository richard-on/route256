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

func (d *Domain) checkAvailable(stocks []Stock, need int64) error {
	for _, s := range stocks {
		need -= int64(s.Count)
		if need <= 0 {
			return nil
		}
	}

	return errors.New("not enough in stock")
}

// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
func (d *Domain) CreateOrder(ctx context.Context, user int64, items []Item) (int64, error) {

	var orderID int64
	orderID, err := d.Repository.CreateOrder(ctx, OrderInfo{
		Status: NewOrder,
		User:   user,
		Items:  items,
	})
	if err != nil {
		return 0, err
	}

	err = d.TransactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) (err error) {
		stocks := make([][]Stock, len(items))
		for i, item := range items {
			stocks[i], err = d.Repository.GetStocks(ctx, item.SKU)
			if err != nil {
				return err
			}
			if err = d.checkAvailable(stocks[i], int64(item.Count)); err != nil {
				changeErr := d.Repository.ChangeOrderStatus(ctx, orderID, Failed)
				if changeErr != nil {
					return changeErr
				}

				return fmt.Errorf("order %v: sku %v: request %v items: %w", orderID, item.SKU, item.Count, err)
			}
		}

		var reserveAmount uint64 = 0
		for i, item := range items {
			toReserve := uint64(item.Count)
			for _, stock := range stocks[i] {
				if toReserve >= stock.Count {
					reserveAmount = stock.Count
				} else {
					reserveAmount = toReserve
				}
				toReserve -= reserveAmount

				if err = d.Repository.DecreaseCount(ctx, stock.WarehouseID, item.SKU, reserveAmount); err != nil {
					return errors.WithMessagef(err, "counting item with sku %v", item.SKU)
				}

				if err = d.Repository.ReserveItem(ctx, stock.WarehouseID, item.SKU, reserveAmount); err != nil {
					return errors.WithMessagef(err, "reserving item with sku %v", item.SKU)
				}

				if toReserve == 0 {
					break
				}
			}
		}

		err = d.Repository.ChangeOrderStatus(ctx, orderID, AwaitingPayment)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
