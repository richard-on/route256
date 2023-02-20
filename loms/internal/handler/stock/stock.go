package stock

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
	SKU uint32 `json:"sku"`
}

type Response struct {
	Stocks []Item `json:"stocks"`
}

func (r Request) Validate() error {
	if r.SKU == 0 {
		return handler.ErrEmptySKU
	}

	return nil
}

type Item struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

func (h *Handler) Handle(ctx context.Context, request Request) (Response, error) {
	log.Printf("stock: %+v", request)

	return Response{
		Stocks: []Item{
			{
				WarehouseID: 12,
				Count:       5,
			},
		},
	}, nil
}
