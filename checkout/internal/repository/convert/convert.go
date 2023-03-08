package convert

import (
	"github.com/jackc/pgtype"
)

func ToSchemaItem(item domain.Item) schema.Item {
	return schema.Item{
		SKU:   int64(item.SKU),
		Count: int32(item.Count),
	}
}

func ToDomainItem(item schema.Item) domain.Item {
	return domain.Item{
		SKU:   uint32(item.SKU),
		Count: uint16(item.Count),
	}
}

func ToSchemaItemArray(items []domain.Item) schema.ItemArray {
	dbItems := make([]schema.Item, 0, len(items))
	for _, item := range items {
		dbItems = append(dbItems, ToSchemaItem(item))
	}

	return schema.ItemArray{
		Elements:   dbItems,
		Dimensions: []pgtype.ArrayDimension{{Length: int32(len(dbItems)), LowerBound: 1}},
		Status:     pgtype.Present,
	}
}

func ToDomainItems(items schema.ItemArray) []domain.Item {
	domainItems := make([]domain.Item, 0, len(items.Elements))
	for _, item := range items.Elements {
		domainItems = append(domainItems, ToDomainItem(item))
	}

	return domainItems
}

func ToDomainOrderInfo(orderInfo schema.OrderInfo) domain.OrderInfo {
	items := ToDomainItems(orderInfo.Items)

	return domain.OrderInfo{
		Status: domain.Status(orderInfo.Status),
		User:   orderInfo.UserID,
		Items:  items,
	}
}

func ToDomainStock(stocks []schema.Stock) []domain.Stock {
	domainStocks := make([]domain.Stock, 0, len(stocks))
	for _, stock := range stocks {
		domainStocks = append(domainStocks, domain.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       uint64(stock.Count),
		})
	}

	return domainStocks
}
