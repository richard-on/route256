package loms

import (
	"context"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/lib/client/wrapper"
)

type StockRequest struct {
	SKU uint32 `json:"sku"`
}

type StockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type StockResponse struct {
	Stock []StockItem `json:"stocks"`
}

func (c *Client) Stock(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	request := StockRequest{SKU: sku}

	response, err := wrapper.NewRequest(ctx, c.urlStock, http.MethodPost, request, StockResponse{})
	if err != nil {
		return nil, err
	}

	stocks := make([]domain.Stock, 0, len(response.Stock))
	for _, s := range response.Stock {
		stocks = append(stocks, domain.Stock{
			WarehouseID: s.WarehouseID,
			Count:       s.Count,
		})
	}

	return stocks, nil
}
