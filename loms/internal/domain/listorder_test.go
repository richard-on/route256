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

func TestListOrder(t *testing.T) {
	t.Parallel()
	type LOMSRepoMockFunc func(mc *minimock.Controller) LOMSRepo
	type ConnMockFunc func(mc *minimock.Controller) transactor.Conn

	type args struct {
		ctx     context.Context
		orderID int64
	}

	var (
		mc      = minimock.NewController(t)
		tx      = txMocks.NewTxMock(t)
		ctx     = context.Background()
		ctxTx   = context.WithValue(ctx, transactor.Key, tx)
		orderID = int64(gofakeit.Number(1, 1<<31))

		items = []model.Item{
			{
				SKU:   gofakeit.Uint32(),
				Count: gofakeit.Uint16(),
			},
		}

		orderInfo = model.Order{
			Status: model.AwaitingPayment,
			User:   int64(gofakeit.Number(1, 1<<31)),
			Items:  items,
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         model.Order
		err          error
		lomsRepoMock LOMSRepoMockFunc
		connMockFunc ConnMockFunc
	}{
		{
			name: "Positive",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: orderInfo,
			err:  nil,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.ListOrderInfoMock.Expect(ctxTx, orderID).Return(orderInfo, nil)
				mock.ListOrderItemsMock.Expect(ctxTx, orderID).Return(items, nil)

				return mock
			},
			connMockFunc: func(mc *minimock.Controller) transactor.Conn {
				mock := txMocks.NewConnMock(mc)
				mock.BeginTxMock.Expect(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Return(tx, nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrNoOrder",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: model.Order{},
			err:  ErrEmptyOrder,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.ListOrderInfoMock.Expect(ctxTx, orderID).Return(model.Order{}, ErrEmptyOrder)

				return mock
			},
			connMockFunc: func(mc *minimock.Controller) transactor.Conn {
				mock := txMocks.NewConnMock(mc)
				mock.BeginTxMock.Expect(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Return(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
		{
			name: "ErrNoOrderItems",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: model.Order{},
			err:  ErrNoOrderItems,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.ListOrderInfoMock.Expect(ctxTx, orderID).Return(orderInfo, nil)
				mock.ListOrderItemsMock.Expect(ctxTx, orderID).Return(nil, ErrNoOrderItems)

				return mock
			},
			connMockFunc: func(mc *minimock.Controller) transactor.Conn {
				mock := txMocks.NewConnMock(mc)
				mock.BeginTxMock.Expect(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Return(tx, nil)
				tx.RollbackMock.Expect(ctx).Return(nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(tt.lomsRepoMock(mc), transactor.New(tt.connMockFunc(mc)))

			order, err := domain.ListOrder(tt.args.ctx, tt.args.orderID)
			require.Equal(t, tt.want, order)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
