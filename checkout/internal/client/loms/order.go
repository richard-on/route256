package loms

import (
	"context"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/lib/client/wrapper"
)

type OrderRequest struct {
	User int64 `json:"user"`
}

type OrderResponse struct {
	OrderID int64 `json:"orderID"`
}

func (c *Client) Order(ctx context.Context, user int64) (domain.OrderInfo, error) {
	request := OrderRequest{User: user}

	resp, err := wrapper.NewRequest(ctx, c.urlOrder, http.MethodPost, request, OrderResponse{})
	if err != nil {
		return domain.OrderInfo{}, err
	}

	return domain.OrderInfo{
		OrderID: resp.OrderID,
	}, nil
}
