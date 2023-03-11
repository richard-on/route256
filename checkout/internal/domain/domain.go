package domain

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
)

// StockChecker is the interface used for obtaining Stock of a product in all warehouses.
type StockChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]*model.Stock, error)
}

// ProductLister is the interface used for obtaining ProductInfo.
type ProductLister interface {
	GetProduct(ctx context.Context, sku uint32) (model.ProductInfo, error)
}

// OrderCreator is the interface used for order creation and obtaining Order Info.
type OrderCreator interface {
	CreateOrder(ctx context.Context, user int64, items []model.Item) (int64, error)
}

// Transactor is the interface that provides abstraction for different transaction isolation levels.
type Transactor interface {

	// RunReadCommitted runs DB operations provided to it as a transaction with read committed isolation level.
	//
	// Note: You should always use ctxTX Context inside transaction block.
	// Do not use the context passed as the first parameter.
	RunReadCommitted(ctx context.Context, f func(ctxTX context.Context) error) error

	// RunRepeatableRead runs DB operations provided to it as a transaction with repeatable read isolation level.
	//
	// Note: You should always use ctxTX Context inside transaction block.
	// Do not use the context passed as the first parameter.
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error

	// RunSerializable runs DB operations provided to it as a transaction with serializable isolation level.
	//
	// Note: You should always use ctxTX Context inside transaction block.
	// Do not use the context passed as the first parameter.
	RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error
}

type CheckoutRepo interface {
	AddToCart(ctx context.Context, userID int64, modelItem model.Item) error

	GetItemCartCount(ctx context.Context, userID int64, modelItem model.Item) (int32, error)
	DecreaseItemCartCount(ctx context.Context, userID int64, modelItem model.Item) error
	DeleteItemCart(ctx context.Context, userID int64, sku uint32) error
	ClearCart(ctx context.Context, userID int64) error

	GetCartItems(ctx context.Context, userID int64) ([]model.Item, error)
}

// Domain represents business logic of checkout service. It wraps interfaces used in a service.
type Domain struct {
	CheckoutRepo
	Transactor
	StockChecker
	ProductLister
	OrderCreator
}

// New creates a new Domain.
func New(repo CheckoutRepo, transactor Transactor,
	checker StockChecker, lister ProductLister, creator OrderCreator) *Domain {
	return &Domain{
		repo,
		transactor,
		checker,
		lister,
		creator,
	}
}
