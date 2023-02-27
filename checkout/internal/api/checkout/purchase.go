package checkout

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
)

// Purchase creates an order containing all products in a user's cart.
func (c *Checkout) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, err
	}

	orderInfo, err := c.domain.CreateOrder(ctx, req.GetUser())
	if err != nil {
		return nil, err
	}

	resp := checkout.PurchaseResponse{
		OrderId: orderInfo.OrderID,
	}

	return &resp, nil
}
