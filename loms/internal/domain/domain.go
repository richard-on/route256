package domain

import "context"

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type Repository interface {
	InsertOrderInfo(ctx context.Context, order OrderInfo) (int64, error)
	InsertOrderItems(ctx context.Context, orderID int64, domainItems []Item) error
	ChangeOrderStatus(ctx context.Context, orderID int64, status Status) error

	GetStocks(ctx context.Context, sku uint32) ([]Stock, error)
	DecreaseStock(ctx context.Context, sku int64, stock Stock) error
	IncreaseStock(ctx context.Context, sku int64, stock Stock) error
	ReserveItem(ctx context.Context, orderID int64, sku int64, stock Stock) error
	//RemoveItemsFromReserved(ctx context.Context, orderID int64) error
	RemoveItemsFromReserved(ctx context.Context, orderID int64) ([]int64, []Stock, error)

	ListOrderInfo(ctx context.Context, orderID int64) (OrderInfo, error)
	ListOrderItems(ctx context.Context, orderID int64) ([]Item, error)

	PayOrder(ctx context.Context, orderID int64) error
	CancelOrder(ctx context.Context, orderID int64) error
}

// Domain represents business logic of Logistics and Order Management System.
// It should wrap interfaces used in a service.
type Domain struct {
	Repository
	TransactionManager
}

// New creates a new Domain.
func New(repo Repository, txManager TransactionManager) *Domain {
	return &Domain{
		Repository:         repo,
		TransactionManager: txManager,
	}
}
