package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

func (i *Implementation) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {

	stocks, err := i.domain.Stocks(ctx, req.GetSku())
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
