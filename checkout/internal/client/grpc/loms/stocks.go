package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// Stocks calls loms.Stocks to check product availability using LOMS gRPC client.
func (c *Client) Stocks(ctx context.Context, sku uint32) ([]*model.Stock, error) {
	resp, err := c.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: sku})
	if err != nil {
		return nil, err
	}

	stocksResp := resp.GetStocks()

	stocks := make([]*model.Stock, 0, len(stocksResp))
	for _, item := range stocksResp {
		stocks = append(stocks, &model.Stock{
			WarehouseID: item.GetWarehouseId(),
			Count:       item.GetCount(),
		})
	}

	return stocks, nil
}
