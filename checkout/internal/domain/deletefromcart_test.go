package domain

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain/mocks"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/transactor"
	txMocks "gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/repository/transactor/mocks"
)

func TestDeleteFromCart(t *testing.T) {
	t.Parallel()
	type checkoutRepoMockFunc func(mc *minimock.Controller) CheckoutRepo
	type connMockFunc func(mc *minimock.Controller) transactor.Conn

	type args struct {
		ctx  context.Context
		user int64
		item model.Item
	}

	var (
		mc     = minimock.NewController(t)
		tx     = mocks.NewTxMock(t)
		ctx    = context.Background()
		ctxTx  = context.WithValue(ctx, transactor.Key, tx)
		userID = int64(gofakeit.Number(1, 100000000))
		sku    = gofakeit.Uint32()
		count  = gofakeit.Number(2, 5)

		item = model.Item{
			SKU:   sku,
			Count: uint16(count),
		}

		errDelete   = errors.New("error when deleting")
		errDecrease = errors.New("error when decreasing count")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name             string
		args             args
		err              error
		checkoutRepoMock checkoutRepoMockFunc
		connMockFunc     connMockFunc
	}{
		{
			name: "Positive-RecordDeleted",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: nil,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctxTx, userID, item).Return(int32(item.Count), nil)
				mock.DeleteItemCartMock.Expect(ctxTx, userID, sku).Return(nil)

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
			name: "Positive-RecordDecreased",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: nil,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctxTx, userID, item).
					Return(int32(item.Count+1), nil)
				mock.DecreaseItemCartCountMock.Expect(ctxTx, userID, item).
					Return(nil)

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
			name: "ErrNotEnoughInCart",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: ErrNotEnoughInCart,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctxTx, userID, item).
					Return(int32(item.Count-1), nil)

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
			name: "ErrNotInCart",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: ErrNotInCart,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctxTx, userID, item).Return(0, ErrNotInCart)

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
			name: "ErrDeleteFailed",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: errDelete,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctxTx, userID, item).Return(int32(item.Count), nil)
				mock.DeleteItemCartMock.Expect(ctxTx, userID, sku).Return(errDelete)

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
			name: "ErrDecreaseFailed",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: errDecrease,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctxTx, userID, item).Return(int32(item.Count+1), nil)
				mock.DecreaseItemCartCountMock.Expect(ctxTx, userID, item).Return(errDecrease)

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

			domain := NewMockDomain(tt.checkoutRepoMock(mc), transactor.New(tt.connMockFunc(mc)))

			err := domain.DeleteFromCart(tt.args.ctx, tt.args.user, tt.args.item)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
