// Package domain provides business-logic for Logistics and Order Management System.
package domain

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
)

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

// LOMSRepo is the interface that provides methods used in LOMS Repository layer.
type LOMSRepo interface {
	InsertOrderInfo(ctx context.Context, order model.Order) (int64, error)
	InsertOrderItems(ctx context.Context, orderID int64, domainItems []model.Item) error
	ChangeOrderStatus(ctx context.Context, orderID int64, status model.Status) error

	GetStocks(ctx context.Context, sku uint32) ([]model.Stock, error)
	DecreaseStock(ctx context.Context, sku int64, stock model.Stock) error
	IncreaseStock(ctx context.Context, sku int64, stock model.Stock) error
	ReserveItem(ctx context.Context, orderID int64, sku int64, stock model.Stock) error
	RemoveItemsFromReserved(ctx context.Context, orderID int64) ([]int64, []model.Stock, error)

	ListOrderInfo(ctx context.Context, orderID int64) (model.Order, error)
	ListOrderItems(ctx context.Context, orderID int64) ([]model.Item, error)

	PayOrder(ctx context.Context, orderID int64) error
	CancelOrder(ctx context.Context, orderID int64) error
}

// Domain represents business logic of Logistics and Order Management System.
// It should wrap interfaces used in a service.
type Domain struct {
	LOMSRepo
	Transactor
}

// New creates a new Domain.
func New(repo LOMSRepo, tx Transactor) *Domain {
	return &Domain{
		LOMSRepo:   repo,
		Transactor: tx,
	}
}
