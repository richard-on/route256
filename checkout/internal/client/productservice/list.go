package productservice

import (
	"context"
	"net/http"
	"route256/lib/client/wrapper"
)

type ListSKURequest struct {
	Token      string `json:"token"`
	StartAfter uint32 `json:"startAfterSku"`
	Count      uint32 `json:"count"`
}

type ListSKUResponse struct {
	SKUList []uint32 `json:"skus"`
}

func (c *Client) ListSKU(ctx context.Context, startAfterSku, count uint32) ([]uint32, error) {
	request := ListSKURequest{
		Token:      c.token,
		StartAfter: startAfterSku,
		Count:      count,
	}

	response, err := wrapper.NewRequest(ctx, c.urlList, http.MethodPost, request, ListSKUResponse{})
	if err != nil {
		return nil, err
	}

	skus := make([]uint32, 0, len(response.SKUList))
	skus = append(skus, response.SKUList...)

	return skus, nil
}
