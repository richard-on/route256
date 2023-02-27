package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc"
)

// Client represents LOMS methods that can be called from this client.
type Client interface {
	CreateOrder(ctx context.Context, user int64) (domain.OrderInfo, error)
	Stocks(ctx context.Context, sku uint32) ([]*domain.Stock, error)
}

// Client is a wrapper for LOMS gRPC client.
type client struct {
	lomsClient loms.LOMSClient
}

// NewClient creates new LOMS gRPC client.
func NewClient(cc *grpc.ClientConn) Client {
	return &client{
		lomsClient: loms.NewLOMSClient(cc),
	}
}
