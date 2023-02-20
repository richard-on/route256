package productservice

import (
	"context"
	"net/http"
	"route256/checkout/internal/domain/product"
	"route256/lib/client/wrapper"
)

type ProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type ProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, sku uint32) (product.Info, error) {
	request := ProductRequest{
		Token: c.token,
		SKU:   sku,
	}

	response, err := wrapper.NewRequest(ctx, c.urlList, http.MethodPost, request, ProductResponse{})
	if err != nil {
		return product.Info{}, err
	}

	return product.Info{
		Name:  response.Name,
		Price: response.Price,
	}, nil
}
