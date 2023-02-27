package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CancelOrder cancels order, makes previously reserved products available.
func (l *LOMS) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error) {
	err := validateOrder(req.GetOrderId())
	if err != nil {
		return nil, err
	}

	err = l.domain.CancelOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
