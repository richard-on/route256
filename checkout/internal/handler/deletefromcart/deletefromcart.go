package deletefromcart

import (
	"context"
	"log"
	"route256/checkout/internal/domain/cart"
	"route256/checkout/internal/handler"
)

type Handler struct {
	cart *cart.Delete
}

func New() *Handler {
	return &Handler{
		cart: cart.NewDeleter(),
	}
}

type Request struct {
	User  int64  `json:"user"`
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (r Request) Validate() error {
	if r.User == 0 {
		return handler.ErrEmptyUser
	}

	if r.Sku == 0 {
		return handler.ErrEmptySKU
	}

	if r.Count == 0 {
		return handler.ErrZeroCount
	}

	return nil
}

type Response struct{}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("deleteFromCart: %+v", req)

	var response Response

	err := h.cart.Delete(ctx, req.User, req.Sku, req.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}
