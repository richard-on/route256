// Package convert converts objects between Database and business logic implementations.
package convert

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/schema"
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

// ToModelItemSlice converts []schema.Item to []model.Item.
func ToModelItemSlice(items []schema.Item) []model.Item {
	modelItems := make([]model.Item, 0, len(items))
	for _, item := range items {
		modelItems = append(modelItems, ToModelItem(item))
	}

	return modelItems
}
