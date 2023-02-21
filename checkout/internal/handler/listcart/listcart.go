package listcart

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handler"

	"github.com/rs/zerolog/log"
)

type Handler struct {
	domain *domain.Domain
}

func New(domain *domain.Domain) *Handler {
	return &Handler{
		domain: domain,
	}
}

type Request struct {
	User int64 `json:"user"`
}

func (r Request) Validate() error {
	if r.User <= 0 {
		return handler.ErrEmptyUser
	}

	return nil
}

type Response struct {
	Items      []domain.Item `json:"items"`
	TotalPrice uint32        `json:"totalPrice"`
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("listCart: %+v", req)

	items, err := h.domain.ListCart(ctx, req.User)
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
