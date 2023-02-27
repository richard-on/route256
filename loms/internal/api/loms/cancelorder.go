package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error) {

	err := i.domain.CancelOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
