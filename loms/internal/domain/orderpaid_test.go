package domain

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain/mocks"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor"
	txMocks "gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/repository/transactor/mocks"
	"testing"
)

func TestOrderPaid(t *testing.T) {
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

		stocks = []model.Stock{
			{
				WarehouseID: int64(gofakeit.Number(1, 100)),
				Count:       uint64(gofakeit.Number(1, 100000)),
			},
		}

		payOrderErr = errors.New("pay order error")
		reserveErr  = errors.New("items remove error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         []model.Stock
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
			want: stocks,
			err:  nil,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.PayOrderMock.Expect(ctxTx, orderID).Return(nil)
				mock.RemoveItemsFromReservedMock.Expect(ctxTx, orderID).Return(nil, nil, nil)
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
			name: "ErrPayOrder",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: nil,
			err:  payOrderErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.PayOrderMock.Expect(ctxTx, orderID).Return(payOrderErr)
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
			name: "ErrRemoveItems",
			args: args{
				ctx:     ctx,
				orderID: orderID,
			},
			want: nil,
			err:  reserveErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.PayOrderMock.Expect(ctxTx, orderID).Return(nil)
				mock.RemoveItemsFromReservedMock.Expect(ctxTx, orderID).Return(nil, nil, reserveErr)
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

			err := domain.OrderPaid(tt.args.ctx, tt.args.orderID)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
