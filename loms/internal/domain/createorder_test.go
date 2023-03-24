package domain

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain/mocks"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor"
	txMocks "gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor/mocks"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()
	type LOMSRepoMockFunc func(mc *minimock.Controller) LOMSRepo
	type ConnMockFunc func(mc *minimock.Controller) transactor.Conn

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
		tx      = txMocks.NewTxMock(t)
		ctx     = context.Background()
		ctxTx   = context.WithValue(ctx, transactor.Key, tx)
		orderID = int64(gofakeit.Number(1, 1<<31))
		user    = int64(gofakeit.Number(1, 1<<31))

		itemNum      = 5
		itemOrderNum = 3
		items        []model.Item
		itemStocks   []itemStock
	)

	t.Cleanup(mc.Finish)

	for i := 0; i < itemNum; i++ {
		sku := gofakeit.Uint32()
		count := uint16(gofakeit.Number(1, 10))

		items = append(items, model.Item{
			SKU:   sku,
			Count: count,
		})

		stockNum := gofakeit.Number(1, 100)
		stocks := make([]model.Stock, 0, stockNum)
		for j := 0; j < gofakeit.Number(1, 100); j++ {
			stocks = append(stocks, model.Stock{
				WarehouseID: int64(gofakeit.Number(1, 100)),
				Count:       uint64(gofakeit.Number(10, 100000)),
			})
		}
		itemStocks = append(itemStocks, itemStock{
			item: model.Item{
				SKU:   sku,
				Count: count,
			},
			stocks: stocks,
		})
	}

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		lomsRepoMock LOMSRepoMockFunc
		connMockFunc ConnMockFunc
	}{
		{
			name: "Positive",
			args: args{
				ctx:   ctx,
				user:  user,
				items: items[:itemOrderNum],
			},
			want: orderID,
			err:  nil,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.InsertOrderInfoMock.Expect(ctxTx, model.Order{Status: model.NewOrder, User: user}).
					Return(orderID, nil)
				mock.InsertOrderItemsMock.Expect(ctxTx, orderID, items[:itemOrderNum]).Return(nil)

				for i := 0; i < itemOrderNum; i++ {
					mock.GetStocksMock.When(ctxTx, items[i].SKU).Then(itemStocks[i].stocks, nil)

					// In this scenario, we assume that the first stock is enough for the order.
					itemStocks[i].stocks[0].Count = uint64(items[i].Count)
					mock.DecreaseStockMock.When(ctxTx, int64(items[i].SKU), itemStocks[i].stocks[0]).
						Then(nil)
					mock.ReserveItemMock.When(ctxTx, orderID, int64(items[i].SKU), itemStocks[i].stocks[0]).
						Then(nil)
				}

				mock.ChangeOrderStatusMock.Expect(ctxTx, orderID, model.AwaitingPayment).Return(nil)

				return mock
			},
			connMockFunc: func(mc *minimock.Controller) transactor.Conn {
				mock := txMocks.NewConnMock(mc)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Then(tx, nil)
				mock.BeginTxMock.When(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}).Then(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(tt.lomsRepoMock(mc), transactor.New(tt.connMockFunc(mc)))

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
