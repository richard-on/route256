package productservice

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	product "gitlab.ozon.dev/rragusskiy/homework-1/productservice/pkg/checkout"
)

// GetProduct calls productService.GetProduct for a given product and returns domain.ProductInfo.
func (c *Client) GetProduct(ctx context.Context, sku uint32) (model.ProductInfo, error) {
	c.rateLimit.Take()

	resp, err := c.productClient.GetProduct(ctx, &product.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	})
	if err != nil {
		return model.ProductInfo{}, err
	}

	return model.ProductInfo{
		Name:  resp.GetName(),
		Price: resp.GetPrice(),
	}, nil
}
