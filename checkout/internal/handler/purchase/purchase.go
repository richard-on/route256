package purchase

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

type Response struct {
	OrderID int64 `json:"orderID"`
}

func (r Request) Validate() error {
	if r.User <= 0 {
		return handler.ErrEmptyUser
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("listCart: %+v", req)

	orderInfo, err := h.domain.CreateOrder(ctx, req.User)
	if err != nil {
		return Response{}, err
	}

	return Response{
		OrderID: orderInfo.OrderID,
	}, nil
}
