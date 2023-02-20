package loms

import (
	"context"
	"net/http"
	"route256/lib/client/wrapper"
)

type OrderRequest struct {
	User int64 `json:"user"`
}

type OrderResponse struct {
}

func (c *Client) Order(ctx context.Context, user int64) error {
	request := OrderRequest{User: user}

	_, err := wrapper.NewRequest(ctx, c.urlOrder, http.MethodPost, request, OrderResponse{})
	if err != nil {
		return err
	}

	return nil
}
