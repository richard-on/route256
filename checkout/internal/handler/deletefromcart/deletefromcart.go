package deletefromcart

import (
	"context"
	"route256/checkout/internal/handler"

	"github.com/rs/zerolog/log"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type Request struct {
	User  int64  `json:"user"`
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (r Request) Validate() error {
	if r.User <= 0 {
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

	return Response{}, nil
}
