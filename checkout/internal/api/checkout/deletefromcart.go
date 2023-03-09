package checkout

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteFromCart deletes a product from a user's cart.
func (c *Checkout) DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
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

	err = c.domain.DeleteFromCart(ctx, req.GetUser(), model.Item{
		SKU:   req.GetSku(),
		Count: uint16(req.GetCount()),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
