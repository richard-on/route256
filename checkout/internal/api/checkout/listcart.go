package checkout

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/pkg/checkout"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListCart lists all products that are currently in a user's cart.
func (c *Checkout) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {
	err := validateUser(req.GetUser())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	items, totalPrice, err := c.domain.ListCart(ctx, req.GetUser())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyCart):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case status.Code(err) == codes.ResourceExhausted:
			return nil, err
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	itemsResp := make([]*checkout.Item, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, &checkout.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
			Name:  item.ProductInfo.Name,
			Price: item.ProductInfo.Price,
		})
	}

	resp := checkout.ListCartResponse{
		Items:      itemsResp,
		TotalPrice: totalPrice,
	}

	return &resp, nil
}
