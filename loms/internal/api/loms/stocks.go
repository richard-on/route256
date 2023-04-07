package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Stocks returns a number of available products with a given SKU in different warehouses.
func (l *LOMS) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {
	err := validateSKU(req.GetSku())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stocks, err := l.domain.Stocks(ctx, req.GetSku())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &loms.StocksResponse{
		Stocks: convert.ToProtoStockSlice(stocks),
	}, nil
}
