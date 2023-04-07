package checkout

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
)

// Purchase creates an order containing all products in a user's cart.
func (c *Checkout) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orderID, err := c.domain.CreateOrder(ctx, req.GetUser())
	switch {
	case errors.Is(err, domain.ErrEmptyCart):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case status.Code(err) == codes.FailedPrecondition:
		return nil, err
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := checkout.PurchaseResponse{
		OrderId: orderID,
	}

	return &resp, nil
}
