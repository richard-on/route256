package create

import (
	"context"
	"route256/loms/internal/handler"

	"github.com/rs/zerolog/log"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type Request struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (r Request) Validate() error {
	if r.User <= 0 {
		return handler.ErrEmptyUser
	}

	for _, item := range r.Items {
		if item.SKU == 0 {
			return handler.ErrEmptySKU
		}
		if item.Count == 0 {
			return handler.ErrZeroCount
		}
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("createOrder: %+v", request)

	return Response{
		OrderID: 42,
	}, nil
}
