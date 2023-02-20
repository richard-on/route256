package cancel

import (
	"context"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/handler"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type Request struct {
	OrderID int64 `json:"orderID"`
}

type Response struct{}

func (r Request) Validate() error {
	if r.OrderID <= 0 {
		return handler.ErrInvalidOrder
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("cancelOrder: %+v", request)

	return Response{}, nil
}
