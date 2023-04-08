package checkout

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteFromCart deletes a product from a user's cart.
func (c *Checkout) DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = validateSKU(req.GetSku())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = validateCount(req.GetCount())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.domain.DeleteFromCart(ctx, req.GetUser(), model.Item{
		SKU:   req.GetSku(),
		Count: uint16(req.GetCount()),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotEnoughInCart):
			return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, domain.ErrNotInCart):
			return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
