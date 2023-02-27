package checkout

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AddToCart adds a product to a user's cart.
func (c *Checkout) AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, err
	}
	err = validateSKU(req.GetSku())
	if err != nil {
		return nil, err
	}
	err = validateCount(req.GetCount())
	if err != nil {
		return nil, err
	}

	err = c.domain.AddToCart(ctx, req.GetUser(), req.GetSku(), uint16(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
