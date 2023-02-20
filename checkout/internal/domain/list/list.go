package list

import (
	"context"
	"route256/checkout/internal/domain/product"
)

type Lister interface {
	GetProduct(ctx context.Context, sku uint32) (product.Info, error)
}

type List struct {
	lister Lister
}

func New(lister Lister) *List {
	return &List{
		lister: lister,
	}
}
