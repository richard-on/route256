package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// Stocks returns a number of available products with a given SKU in different warehouses.
func (l *LOMS) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {

	stocks, err := l.domain.Stocks(ctx, req.GetSku())
	if err != nil {
		return nil, err
	}

	stocksLOMS := make([]*loms.Stock, 0, len(stocks))
	for _, stock := range stocks {
		stocksLOMS = append(stocksLOMS, &loms.Stock{
			WarehouseId: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return &loms.StocksResponse{
		Stocks: stocksLOMS,
	}, nil
}
