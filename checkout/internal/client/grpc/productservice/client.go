package productservice

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/lib/ratelimit"
	product "gitlab.ozon.dev/rragusskiy/homework-1/productservice/pkg/checkout"
	"google.golang.org/grpc"
)

// Client is a wrapper for product service gRPC client.
type Client struct {
	productClient product.ProductServiceClient
	token         string
	rateLimit     ratelimit.Limiter
}

// NewClient creates new product service LOMS gRPC client.
func NewClient(cc *grpc.ClientConn, token string, rateLimit ratelimit.Limiter) *Client {
	return &Client{
		productClient: product.NewProductServiceClient(cc),
		token:         token,
		rateLimit:     rateLimit,
	}
}
