package purchase

import (
	"context"
	"log"
	"route256/checkout/internal/domain/purchase"
	"route256/checkout/internal/handler"
)

type Handler struct {
	order *purchase.Order
}

func New(order *purchase.Order) *Handler {
	return &Handler{
		order: order,
	}
}

type Request struct {
	User int64 `json:"user"`
}

func (r Request) Validate() error {
	if r.User == 0 {
		return handler.ErrEmptyUser
	}

	return nil
}

type Response struct{}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("listCart: %+v", req)

	err := h.order.Create(ctx, req.User)
	if err != nil {
		return Response{}, err
	}

	return Response{}, nil
}
