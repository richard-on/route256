package productservice

import (
	product "gitlab.ozon.dev/rragusskiy/homework-1/productservice/pkg/checkout"
	"google.golang.org/grpc"
)

// Client is a wrapper for product service gRPC client.
type Client struct {
	productClient product.ProductServiceClient
	token         string
}

// NewClient creates new product service LOMS gRPC client.
func NewClient(cc *grpc.ClientConn, token string) *Client {
	return &Client{
		productClient: product.NewProductServiceClient(cc),
		token:         token,
	}
}
