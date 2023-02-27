package checkout

import (
	"context"
	"math"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	// ErrZeroCount is the error returned when count is zero.
	ErrZeroCount = errors.New("zero count")
	// ErrTypeOverflow is the error returned when value is more than max value for uint16.
	ErrTypeOverflow = errors.New("uint16 overflow")
)

// validateCount validates value of count.
//
// Count should be non-zero and not bigger than 65535 (1<<16 - 1).
func validateCount(count uint32) error {
	if count == 0 {
		return ErrZeroCount
	}
	if count > math.MaxUint16 {
		return ErrTypeOverflow
	}

	return nil
}

// AddToCart adds a product to a user's cart.
func (c *Checkout) AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error) {
	err := validateCount(req.GetCount())
	if err != nil {
		return nil, errors.WithMessage(err, "validating request")
	}

	err = c.domain.AddToCart(ctx, req.GetUser(), req.GetSku(), uint16(req.GetCount()))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
