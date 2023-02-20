package listcart

import (
	"context"
	"log"
	"route256/checkout/internal/domain/list"
	"route256/checkout/internal/handler"
)

type Handler struct {
	list *list.List
}

func New(list *list.List) *Handler {
	return &Handler{
		list: list,
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

type Response struct {
	Items      []list.Item `json:"items"`
	TotalPrice uint32      `json:"totalPrice"`
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("listCart: %+v", req)

	items, err := h.list.List(ctx, req.User)
	if err != nil {
		return Response{}, err
	}

	var totalPrice uint32 = 0
	for _, item := range items {
		totalPrice += item.Price
	}

	return Response{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
