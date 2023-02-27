package productservice

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	product "gitlab.ozon.dev/rragusskiy/homework-1/productservice/pkg/checkout"
	"google.golang.org/grpc"
)

// Client represents product service methods that can be called from this client.
type Client interface {
	GetProduct(ctx context.Context, sku uint32) (domain.ProductInfo, error)
	ListSKU(ctx context.Context, startAfterSku, count uint32) ([]uint32, error)
}

// Client is a wrapper for product service gRPC client.
type client struct {
	productClient product.ProductServiceClient
	token         string
}

// NewClient creates new product service LOMS gRPC client.
func NewClient(cc *grpc.ClientConn, token string) Client {
	return &client{
		productClient: product.NewProductServiceClient(cc),
		token:         token,
	}
}
