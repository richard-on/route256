package checkout

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteFromCart deletes a product from a user's cart.
func (c *Checkout) DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
	err := validateCount(req.GetCount())
	if err != nil {
		return nil, errors.WithMessage(err, "validating request")
	}

	err = c.domain.DeleteFromCart(ctx, req.GetUser(), req.GetSku(), uint16(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
