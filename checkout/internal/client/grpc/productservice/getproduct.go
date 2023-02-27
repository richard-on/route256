package productservice

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	product "gitlab.ozon.dev/rragusskiy/homework-1/productservice/pkg/checkout"
)

// GetProduct calls productService.GetProduct for a given product and returns domain.ProductInfo.
func (c *Client) GetProduct(ctx context.Context, sku uint32) (domain.ProductInfo, error) {
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
