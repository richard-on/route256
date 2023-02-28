package domain

import "context"

// StockChecker is the interface used for obtaining Stock of a product in all warehouses.
type StockChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]*Stock, error)
}

// ProductLister is the interface used for obtaining ProductInfo.
type ProductLister interface {
	GetProduct(ctx context.Context, sku uint32) (ProductInfo, error)
}

// OrderCreator is the interface used for order creation and obtaining OrderInfo.
type OrderCreator interface {
	CreateOrder(ctx context.Context, user int64) (OrderInfo, error)
}

// Domain represents business logic of checkout service. It wraps interfaces used in a service.
type Domain struct {
	stockChecker  StockChecker
	productLister ProductLister
	orderCreator  OrderCreator
}

// New creates a new Domain.
func New(checker StockChecker, lister ProductLister, creator OrderCreator) *Domain {
	return &Domain{
		stockChecker:  checker,
		orderCreator:  creator,
		productLister: lister,
	}
}
