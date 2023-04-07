package domain

//go:generate minimock -i StockChecker -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i ProductLister -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i OrderCreator -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i github.com/jackc/pgx/v4.Tx -o ./mocks/tx_minimock.go -n TxMock
//go:generate minimock -i CheckoutRepo -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/config"
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
	config config.Service
	CheckoutRepo
	Transactor
	StockChecker
	ProductLister
	OrderCreator
}

// New creates a new Domain.
func New(config config.Service, repo CheckoutRepo, transactor Transactor,
	checker StockChecker, lister ProductLister, creator OrderCreator) *Domain {
	return &Domain{
		config,
		repo,
		transactor,
		checker,
		lister,
		creator,
	}
}

// NewMockDomain creates a new mock Domain used for testing.
func NewMockDomain(opts ...any) *Domain {
	d := Domain{}

	for _, v := range opts {
		switch s := v.(type) {
		case config.Service:
			d.config = s
		case CheckoutRepo:
			d.CheckoutRepo = s
		case Transactor:
			d.Transactor = s
		case StockChecker:
			d.StockChecker = s
		case ProductLister:
			d.ProductLister = s
		case OrderCreator:
			d.OrderCreator = s
		}
	}

	return &d
}
