// Package convert converts objects between Database and business logic implementations.
package convert

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/schema"
)

// ToSchemaItem converts model.Item to schema.Item.
func ToSchemaItem(item model.Item) schema.Item {
	return schema.Item{
		SKU:   int64(item.SKU),
		Count: int32(item.Count),
	}
}

// ToModelItem converts schema.Item to model.Item.
func ToModelItem(item schema.Item) model.Item {
	return model.Item{
		SKU:   uint32(item.SKU),
		Count: uint16(item.Count),
	}
}

// ToSchemaItemSlice converts []model.Item to []schema.Item.
func ToSchemaItemSlice(items []model.Item) []schema.Item {
	dbItems := make([]schema.Item, 0, len(items))
	for _, item := range items {
		dbItems = append(dbItems, ToSchemaItem(item))
	}

	return dbItems
}

// ToModelItemSlice converts []schema.Item to []model.Item.
func ToModelItemSlice(items []schema.Item) []model.Item {
	modelItems := make([]model.Item, 0, len(items))
	for _, item := range items {
		modelItems = append(modelItems, ToModelItem(item))
	}

	return modelItems
}

// ToModelOrder converts schema.Order to model.Order.
func ToModelOrder(order schema.Order) model.Order {
	var items []model.Item
	if order.Items != nil {
		items = ToModelItemSlice(order.Items)
	}

	return model.Order{
		Status: model.Status(order.Status),
		User:   order.UserID,
		Items:  items,
	}
}

// ToModelStock converts schema.Stock to model.Stock.
func ToModelStock(stock schema.Stock) model.Stock {
	return model.Stock{
		WarehouseID: stock.WarehouseID,
		Count:       uint64(stock.Count),
	}
}

// ToModelStockSlice converts []schema.Stock to []model.Stock.
func ToModelStockSlice(stocks []schema.Stock) []model.Stock {
	modelStocks := make([]model.Stock, 0, len(stocks))
	for _, stock := range stocks {
		modelStocks = append(modelStocks, ToModelStock(stock))
	}

	return modelStocks
}
