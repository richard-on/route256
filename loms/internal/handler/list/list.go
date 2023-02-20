package list

import (
	"context"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/handler"
)

const (
	StatusNew       = "new"
	StatusAwaiting  = "awaiting payment"
	StatusFailed    = "failed"
	StatusPayed     = "payed"
	StatusCancelled = "cancelled"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type Request struct {
	OrderID int64 `json:"orderID"`
}

type Response struct {
	Status string `json:"status"`
	User   int64  `json:"user"`
	Items  []Item `json:"items"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (r Request) Validate() error {
	if r.OrderID <= 0 {
		return handler.ErrInvalidOrder
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("listOrder: %+v", request)

	return Response{
		Status: StatusNew,
		User:   123,
		Items: []Item{
			{
				SKU:   111111,
				Count: 2,
			},
			{
				SKU:   222222,
				Count: 1,
			},
		},
	}, nil
}
