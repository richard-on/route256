package productservice

import (
	"context"

	product "gitlab.ozon.dev/rragusskiy/homework-1/productservice/pkg/checkout"
)

// ListSKU calls productService.ListSkus for a given number of skus and returns a slice of actual skus.
func (c *Client) ListSKU(ctx context.Context, startAfterSku, count uint32) ([]uint32, error) {
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
