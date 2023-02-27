package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *loms.OrderPayedRequest) (*emptypb.Empty, error) {

	err := i.domain.OrderPayed(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
