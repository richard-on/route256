package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/protobuf/types/known/emptypb"
)

// OrderPaid marks order as paid.
func (l *LOMS) OrderPaid(ctx context.Context, req *loms.OrderPaidRequest) (*emptypb.Empty, error) {
	err := validateOrder(req.GetOrderId())
	if err != nil {
		return nil, err
	}

	err = l.domain.OrderPaid(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
