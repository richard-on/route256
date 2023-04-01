// Package domain provides business-logic for Logistics and Order Management System.
package domain

//go:generate minimock -i github.com/jackc/pgx/v4.Tx -o ./mocks/tx_minimock.go -n TxMock
//go:generate minimock -i LOMSRepo -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i OutboxRepo -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i StatusSender -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"time"

	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model/outbox"
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

	ListOrderInfo(ctx context.Context, orderID int64) (model.Order, error)
	ListOrderItems(ctx context.Context, orderID int64) ([]model.Item, error)
	ListUnpaidOrders(ctx context.Context, paymentWait time.Duration) ([]int64, error)

	CancelOrder(ctx context.Context, orderID int64) error
	PayOrder(ctx context.Context, orderID int64) error
	ChangeOrderStatus(ctx context.Context, orderID int64, status model.Status) error

	GetStocks(ctx context.Context, sku uint32) ([]model.Stock, error)
	DecreaseStock(ctx context.Context, sku int64, stock model.Stock) error
	IncreaseStock(ctx context.Context, sku int64, stock model.Stock) error

	ReserveItem(ctx context.Context, orderID int64, sku int64, stock model.Stock) error
	RemoveItemsFromReserved(ctx context.Context, orderID int64) ([]int64, []model.Stock, error)

	OutboxRepo
}

// OutboxRepo is the interface that provides methods used in LOMS Outbox.
type OutboxRepo interface {
	AddMessageWithKey(ctx context.Context, key string, payload []byte) error
	AddMessageWithoutKey(ctx context.Context, payload []byte) error
	UpdateMessageStatus(ctx context.Context, id int64, status outbox.Status) error
	DeleteMessage(ctx context.Context, id int64) error
	ListUnsent(ctx context.Context) ([]outbox.Message, error)
}

// StatusSender is the interface that provides methods used to send a status update message to the broker.
type StatusSender interface {
	SendWithKey(id int64, key string, message []byte)
}

// Domain represents business logic of Logistics and Order Management System.
// It should wrap interfaces used in a service.
type Domain struct {
	config config.Service
	LOMSRepo
	Transactor
	StatusSender
}

// New creates a new Domain.
func New(config config.Service, repo LOMSRepo, tx Transactor, sender StatusSender) *Domain {
	return &Domain{
		config,
		repo,
		tx,
		sender,
	}
}

// NewMockDomain creates a new mock Domain used for testing.
func NewMockDomain(opts ...any) *Domain {
	d := Domain{}

	for _, v := range opts {
		switch s := v.(type) {
		case config.Service:
			d.config = s
		case LOMSRepo:
			d.LOMSRepo = s
		case Transactor:
			d.Transactor = s
		}
	}

	return &d
}
