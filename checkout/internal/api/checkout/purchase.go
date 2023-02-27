package checkout

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
)

func (i *Implementation) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {

	orderInfo, err := i.domain.CreateOrder(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	resp := checkout.PurchaseResponse{
		OrderId: orderInfo.OrderID,
	}

	return &resp, nil
}
