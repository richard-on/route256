package domain

import (
	"context"
	"testing"

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

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	type LOMSRepoMockFunc func(mc *minimock.Controller) LOMSRepo
	type DBMockFunc func(mc *minimock.Controller) transactor.DB

	type args struct {
		ctx   context.Context
		user  int64
		items []model.Item
	}

	type itemStock struct {
		item   model.Item
		stocks []model.Stock
	}

	var (
		mc      = minimock.NewController(t)
		tx      = mocks.NewTxMock(t)
		ctx     = context.Background()
		ctxTx   = context.WithValue(ctx, transactor.Key, tx)
		cfg     = config.Service{MaxPoolWorkers: 5}
		orderID = int64(gofakeit.Number(1, 1<<31))
		user    = int64(gofakeit.Number(1, 1<<31))

		skus = []uint32{gofakeit.Uint32(), gofakeit.Uint32()}

		itemsToBuy = []model.Item{
			{
				SKU:   skus[0],
				Count: 5,
			},
			{
				SKU:   skus[1],
				Count: 8,
			},
		}

		itemStocks = []itemStock{
			{
				item: itemsToBuy[0],
				stocks: []model.Stock{
					{
						WarehouseID: int64(gofakeit.Number(1, 1000)),
						Count:       3,
					},
					{
						WarehouseID: int64(gofakeit.Number(1, 1000)),
						Count:       3,
					},
				},
			},
			{
				item: itemsToBuy[1],
				stocks: []model.Stock{
					{
						WarehouseID: int64(gofakeit.Number(1, 1000)),
						Count:       6,
					},
					{
						WarehouseID: int64(gofakeit.Number(1, 1000)),
						Count:       10,
					},
				},
			},
		}

		insertInfoErr   = errors.New("insert order info error")
		insertItemsErr  = errors.New("insert order items error")
		getStocksErr    = errors.New("get stocks error")
		changeStatusErr = errors.New("change status error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		lomsRepoMock LOMSRepoMockFunc
		DBMockFunc   DBMockFunc
	}{
		{
			name: "Positive-EnoughStocksInFirstWarehouse",
			args: args{
				ctx:   ctx,
				user:  user,
				items: itemsToBuy,
			},
			want: orderID,
			err:  nil,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(orderID, nil)
				mock.InsertOrderItemsMock.Expect(ctxTx, orderID, itemsToBuy).Return(nil)
				mock.AddMessageWithKeyMock.Return(nil)

				for i := 0; i < len(itemsToBuy); i++ {
					mock.GetStocksMock.When(ctxTx, itemsToBuy[i].SKU).Then(itemStocks[i].stocks, nil)

					itemStocks[i].stocks[0].Count = uint64(itemsToBuy[i].Count)
					mock.DecreaseStockMock.When(ctxTx, int64(itemsToBuy[i].SKU), itemStocks[i].stocks[0]).
						Then(nil)
					mock.ReserveItemMock.When(ctxTx, orderID, int64(itemsToBuy[i].SKU), itemStocks[i].stocks[0]).
						Then(nil)
				}

				mock.ChangeOrderStatusMock.Expect(ctxTx, orderID, model.AwaitingPayment).Return(nil)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrInsertOrderInfo",
			args: args{
				ctx:   ctx,
				user:  user,
				items: itemsToBuy,
			},
			want: 0,
			err:  insertInfoErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(0, insertInfoErr)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrInsertOrderItems",
			args: args{
				ctx:   ctx,
				user:  user,
				items: itemsToBuy,
			},
			want: 0,
			err:  insertItemsErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(orderID, nil)
				mock.InsertOrderItemsMock.Expect(ctxTx, orderID, itemsToBuy).Return(insertItemsErr)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrGetStocks",
			args: args{
				ctx:   ctx,
				user:  user,
				items: itemsToBuy,
			},
			want: 0,
			err:  getStocksErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(orderID, nil)
				mock.InsertOrderItemsMock.Expect(ctxTx, orderID, itemsToBuy).Return(nil)
				mock.GetStocksMock.When(ctxTx, itemsToBuy[0].SKU).Then(itemStocks[0].stocks, getStocksErr)

				mock.ChangeOrderStatusMock.Expect(ctx, orderID, model.Failed).Return(nil)
				mock.AddMessageWithKeyMock.Return(nil)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrDecreaseStock",
			args: args{
				ctx:   ctx,
				user:  user,
				items: itemsToBuy,
			},
			want: 0,
			err:  ErrStockNotExists,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(orderID, nil)
				mock.InsertOrderItemsMock.Expect(ctxTx, orderID, itemsToBuy).Return(nil)

				mock.GetStocksMock.When(ctxTx, itemsToBuy[0].SKU).Then(itemStocks[0].stocks, nil)

				itemStocks[0].stocks[0].Count = uint64(itemsToBuy[0].Count)
				mock.DecreaseStockMock.When(ctxTx, int64(itemsToBuy[0].SKU), itemStocks[0].stocks[0]).
					Then(ErrStockNotExists)

				mock.ChangeOrderStatusMock.Expect(ctx, orderID, model.Failed).Return(nil)
				mock.AddMessageWithKeyMock.Return(nil)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrReserveItem",
			args: args{
				ctx:   ctx,
				user:  user,
				items: itemsToBuy,
			},
			want: 0,
			err:  ErrStockNotExists,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(orderID, nil)
				mock.InsertOrderItemsMock.Expect(ctxTx, orderID, itemsToBuy).Return(nil)

				mock.GetStocksMock.When(ctxTx, itemsToBuy[0].SKU).Then(itemStocks[0].stocks, nil)

				itemStocks[0].stocks[0].Count = uint64(itemsToBuy[0].Count)
				mock.DecreaseStockMock.When(ctxTx, int64(itemsToBuy[0].SKU), itemStocks[0].stocks[0]).
					Then(nil)
				mock.ReserveItemMock.When(ctxTx, orderID, int64(itemsToBuy[0].SKU), itemStocks[0].stocks[0]).
					Then(ErrStockNotExists)

				mock.ChangeOrderStatusMock.Expect(ctx, orderID, model.Failed).Return(nil)
				mock.AddMessageWithKeyMock.Return(nil)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrChangeOrderStatusFailed",
			args: args{
				ctx:   ctx,
				user:  user,
				items: itemsToBuy,
			},
			want: 0,
			err:  changeStatusErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(orderID, nil)
				mock.InsertOrderItemsMock.Expect(ctxTx, orderID, itemsToBuy).Return(nil)
				mock.GetStocksMock.When(ctxTx, itemsToBuy[0].SKU).Then(itemStocks[0].stocks, getStocksErr)

				mock.ChangeOrderStatusMock.Expect(ctx, orderID, model.Failed).Return(changeStatusErr)
				mock.AddMessageWithKeyMock.Return(nil)

				return mock
			},
			DBMockFunc: func(mc *minimock.Controller) transactor.DB {
				mock := txMocks.NewDBMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)
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

			order, err := domain.CreateOrder(tt.args.ctx, tt.args.user, tt.args.items)
			require.Equal(t, tt.want, order)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
