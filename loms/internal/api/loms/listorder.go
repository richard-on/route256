package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/api/loms/convert"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
)

// ListOrder lists order information.
func (l *LOMS) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {
	err := validateOrder(req.GetOrderId())
	if err != nil {
		return nil, err
	}

	orderInfo, err := l.domain.ListOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &loms.ListOrderResponse{
		Status: loms.Status(orderInfo.Status),
		User:   orderInfo.User,
		Items:  convert.ToProtoItemSlice(orderInfo.Items),
	}, nil
}
