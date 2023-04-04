package domain

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain/mocks"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor"
	txMocks "gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor/mocks"
)

func TestCancelOrder(t *testing.T) {
	t.Parallel()
	type LOMSRepoMockFunc func(mc *minimock.Controller) LOMSRepo
	type DBMockFunc func(mc *minimock.Controller) transactor.DB

	type args struct {
		ctx     context.Context
		orderID int64
	}

	var (
		mc      = minimock.NewController(t)
		tx      = mocks.NewTxMock(t)
		ctx     = context.Background()
		ctxTx   = context.WithValue(ctx, transactor.Key, tx)
		cfg     = config.Service{MaxPoolWorkers: 5}
		orderID = int64(gofakeit.Number(1, 1<<31))

		skus   = []int64{int64(gofakeit.Uint32()), int64(gofakeit.Uint32())}
		stocks = []model.Stock{
			{
				WarehouseID: int64(gofakeit.Number(1, 1000)),
				Count:       6,
			},
			{
				WarehouseID: int64(gofakeit.Number(1, 1000)),
				Count:       10,
			},
		}

		reserveErr = errors.New("reserve removal error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		err          error
		lomsRepoMock LOMSRepoMockFunc
		DBMockFunc   DBMockFunc
	}{
		{
			name: "Positive",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: nil,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.CancelOrderMock.Expect(ctxTx, orderID).Return(nil)
				mock.AddMessageWithKeyMock.Return(nil)
				mock.RemoveItemsFromReservedMock.Expect(ctxTx, orderID).Return(skus, stocks, nil)

				for i, sku := range skus {
					mock.IncreaseStockMock.When(ctxTx, sku, stocks[i]).Then(nil)
				}

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrOrderCancelled",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: ErrOrderCancelled,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.CancelOrderMock.Expect(ctxTx, orderID).Return(ErrOrderCancelled)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrReserveRemove",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: reserveErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.CancelOrderMock.Expect(ctxTx, orderID).Return(nil)
				mock.RemoveItemsFromReservedMock.Expect(ctxTx, orderID).Return(skus, stocks, reserveErr)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrStockNotExists",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			err: ErrStockNotExists,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.CancelOrderMock.Expect(ctxTx, orderID).Return(nil)
				mock.CancelOrderMock.Expect(ctxTx, orderID).Return(nil)
				mock.RemoveItemsFromReservedMock.Expect(ctxTx, orderID).Return(skus, stocks, nil)

				mock.IncreaseStockMock.When(ctxTx, skus[0], stocks[0]).Then(ErrStockNotExists)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(cfg, tt.lomsRepoMock(mc), transactor.New(tt.DBMockFunc(mc)))

			err := domain.CancelOrder(tt.args.ctx, tt.args.orderID)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}

func TestCancelUnpaidOrders(t *testing.T) {
	t.Parallel()
	type LOMSRepoMockFunc func(mc *minimock.Controller) LOMSRepo
	type DBMockFunc func(mc *minimock.Controller) transactor.DB

	type args struct {
		ctx            context.Context
		paymentTimeout time.Duration
	}

	var (
		mc             = minimock.NewController(t)
		tx             = mocks.NewTxMock(t)
		ctx            = context.Background()
		ctxTx          = context.WithValue(ctx, transactor.Key, tx)
		cfg            = config.Service{MaxPoolWorkers: 5}
		paymentTimeout = 1 * time.Second

		unpaidOrders = []int64{int64(gofakeit.Number(1, 1<<31)), int64(gofakeit.Number(1, 1<<31))}
		skus         = []int64{int64(gofakeit.Uint32()), int64(gofakeit.Uint32())}
		stocks       = []model.Stock{
			{
				WarehouseID: int64(gofakeit.Number(1, 1000)),
				Count:       6,
			},
			{
				WarehouseID: int64(gofakeit.Number(1, 1000)),
				Count:       10,
			},
		}

		listErr = errors.New("list unpaid orders error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		errs         []error
		lomsRepoMock LOMSRepoMockFunc
		DBMockFunc   DBMockFunc
	}{
		{
			name: "Positive",
			args: args{
				ctx:            ctx,
				paymentTimeout: paymentTimeout,
			},
			errs: nil,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.ListUnpaidOrdersMock.Expect(ctx, paymentTimeout).Return(unpaidOrders, nil)
				mock.AddMessageWithKeyMock.Return(nil)
				for i := 0; i < len(unpaidOrders); i++ {
					mock.CancelOrderMock.When(ctxTx, unpaidOrders[i]).Then(nil)
					mock.RemoveItemsFromReservedMock.When(ctxTx, unpaidOrders[i]).Then(skus, stocks, nil)

				}
				for j, sku := range skus {
					mock.IncreaseStockMock.When(ctxTx, sku, stocks[j]).Then(nil)
				}

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrListUnpaidOrders",
			args: args{
				ctx:            ctx,
				paymentTimeout: paymentTimeout,
			},
			errs: []error{listErr},
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.ListUnpaidOrdersMock.Expect(ctx, paymentTimeout).Return(nil, listErr)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				return txMocks.NewDBMock(mc)
			},
		},
		{
			name: "ErrCancelOrder",
			args: args{
				ctx:            ctx,
				paymentTimeout: paymentTimeout,
			},
			errs: []error{ErrOrderCancelled, ErrOrderCancelled},
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.ListUnpaidOrdersMock.Expect(ctx, paymentTimeout).Return(unpaidOrders, nil)
				for i := 0; i < len(unpaidOrders); i++ {
					mock.CancelOrderMock.When(ctxTx, unpaidOrders[i]).Then(ErrOrderCancelled)
				}

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(cfg, tt.lomsRepoMock(mc), transactor.New(tt.DBMockFunc(mc)))

			errs := domain.CancelUnpaidOrders(tt.args.ctx, tt.args.paymentTimeout)
			if tt.errs != nil {
				for i, err := range errs {
					require.ErrorContains(t, err, tt.errs[i].Error())
				}
			} else {
				require.Nil(t, errs)
			}
		})
	}
}
