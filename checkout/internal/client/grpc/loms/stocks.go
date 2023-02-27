package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// Stocks calls loms.Stocks to check product availability using LOMS gRPC client.
func (c *Client) Stocks(ctx context.Context, sku uint32) ([]*domain.Stock, error) {
	resp, err := c.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: sku})
	if err != nil {
		return nil, err
	}

	stocksResp := resp.GetStocks()

	stocks := make([]*domain.Stock, 0, len(stocksResp))
	for _, item := range stocksResp {
		stocks = append(stocks, &domain.Stock{
			WarehouseID: item.GetWarehouseId(),
			Count:       item.GetCount(),
		})
	}

	return stocks, nil
}
