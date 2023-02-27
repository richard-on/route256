package productservice

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	product "gitlab.ozon.dev/rragusskiy/homework-1/productservice/pkg/checkout"
	"google.golang.org/grpc"
)

type Client interface {
	GetProduct(ctx context.Context, sku uint32) (domain.ProductInfo, error)
	ListSKU(ctx context.Context, startAfterSku, count uint32) ([]uint32, error)
}

type client struct {
	productClient product.ProductServiceClient
	token         string
}

func NewClient(cc *grpc.ClientConn, token string) *client {
	return &client{
		productClient: product.NewProductServiceClient(cc),
		token:         token,
	}
}

func (c *client) GetProduct(ctx context.Context, sku uint32) (domain.ProductInfo, error) {
	resp, err := c.productClient.GetProduct(ctx, &product.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	})
	if err != nil {
		return domain.ProductInfo{}, err
	}

	return domain.ProductInfo{
		Name:  resp.GetName(),
		Price: resp.GetPrice(),
	}, nil
}

func (c *client) ListSKU(ctx context.Context, startAfterSku, count uint32) ([]uint32, error) {
	resp, err := c.productClient.ListSkus(ctx, &product.ListSkusRequest{
		Token:         c.token,
		StartAfterSku: startAfterSku,
		Count:         count,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetSkus(), nil
}
