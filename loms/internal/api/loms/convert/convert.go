// Package convert converts objects between Proto and business logic implementations.
package convert

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// ToModelItem converts *loms.Item to model.Item.
//
// Note: This operation involves narrowing conversion of field Count (uint32 in Proto to uint16 in model).
// Therefore, field Count must be validated before conversion. Failure to do so may produce unexpected results.
func ToModelItem(item *loms.Item) model.Item {
	return model.Item{
		SKU:   item.Sku,
		Count: uint16(item.Count),
	}
}

// ToProtoItem converts model.Item to *loms.Item.
func ToProtoItem(item model.Item) *loms.Item {
	return &loms.Item{
		Sku:   item.SKU,
		Count: uint32(item.Count),
	}
}

// ToProtoItemSlice converts []model.Item to []*loms.Item.
func ToProtoItemSlice(items []model.Item) []*loms.Item {
	protoItems := make([]*loms.Item, 0, len(items))
	for _, item := range items {
		protoItems = append(protoItems, ToProtoItem(item))
	}

	return protoItems
}

// ToProtoStock converts model.Stock to *loms.Stock.
func ToProtoStock(stock model.Stock) *loms.Stock {
	return &loms.Stock{
		WarehouseId: stock.WarehouseID,
		Count:       stock.Count,
	}
}

// ToProtoStockSlice converts []model.Stock to []*loms.Stock.
func ToProtoStockSlice(stocks []model.Stock) []*loms.Stock {
	protoStocks := make([]*loms.Stock, 0, len(stocks))
	for _, stock := range stocks {
		protoStocks = append(protoStocks, ToProtoStock(stock))
	}

	return protoStocks
}
