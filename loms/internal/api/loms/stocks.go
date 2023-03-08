package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// Stocks returns a number of available products with a given SKU in different warehouses.
func (l *LOMS) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {
	err := validateSKU(req.GetSku())
	if err != nil {
		return nil, err
	}

	stocks, err := l.domain.Stocks(ctx, req.GetSku())
	if err != nil {
		return nil, err
	}

	return &loms.StocksResponse{
		Stocks: convert.ToProtoStockSlice(stocks),
	}, nil
}
