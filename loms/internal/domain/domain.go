package domain

import "context"

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type Repository interface {
	CreateOrder(ctx context.Context, order OrderInfo) (int64, error)
	ChangeOrderStatus(ctx context.Context, orderID int64, status Status) error
	ReserveItem(ctx context.Context, warehouseID int64, sku uint32, reserveAmount uint64) error
	DecreaseCount(ctx context.Context, warehouseID int64, sku uint32, count uint64) error

	ListOrder(ctx context.Context, orderID int64) (OrderInfo, error)
	PayOrder(ctx context.Context, orderID int64) error
	CancelOrder(ctx context.Context, orderID int64) error
	GetStocks(ctx context.Context, sku uint32) ([]Stock, error)
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
