package domain

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/config"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain/mocks"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
)

func TestPurchase(t *testing.T) {
	t.Parallel()
	type checkoutRepoMockFunc func(mc *minimock.Controller) CheckoutRepo
	type orderCreatorMockFunc func(mc *minimock.Controller) OrderCreator

	type args struct {
		ctx    context.Context
		userID int64
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		cfg     = config.Service{MaxPoolWorkers: 5}
		itemNum = 5

		orderID   int64 = 1
		userID    int64 = 1
		cartItems []model.Item

		errCreateOrder = errors.New("error creating order")
		errClearCart   = errors.New("error clearing the cart")
	)

	t.Cleanup(mc.Finish)

	for i := 0; i < itemNum; i++ {
		sku := uint32(gofakeit.Number(50000, 5000000))
		count := uint32(gofakeit.Number(1, 100))

		domainItem := model.Item{
			SKU:   sku,
			Count: uint16(count),
		}
		cartItems = append(cartItems, domainItem)
	}

	tests := []struct {
		name             string
		args             args
		want             int64
		err              error
		checkoutRepoMock checkoutRepoMockFunc
		orderCreatorMock orderCreatorMockFunc
	}{
		{
			name: "Positive",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want: orderID,
			err:  nil,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(cartItems, nil)
				mock.ClearCartMock.Expect(ctx, userID).Return(nil)
				return mock
			},
			orderCreatorMock: func(mc *minimock.Controller) OrderCreator {
				mock := mocks.NewOrderCreatorMock(mc)
				mock.CreateOrderMock.Expect(ctx, userID, cartItems).Return(orderID, nil)
				return mock
			},
		},
		{
			name: "ErrEmptyCart-GetCartItems",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want: 0,
			err:  ErrEmptyCart,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(nil, ErrEmptyCart)
				return mock
			},
			orderCreatorMock: func(mc *minimock.Controller) OrderCreator {
				return mocks.NewOrderCreatorMock(mc)
			},
		},
		{
			name: "ErrCreateOrder",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want: 0,
			err:  errCreateOrder,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(cartItems, nil)
				return mock
			},
			orderCreatorMock: func(mc *minimock.Controller) OrderCreator {
				mock := mocks.NewOrderCreatorMock(mc)
				mock.CreateOrderMock.Expect(ctx, userID, cartItems).Return(0, errCreateOrder)
				return mock
			},
		},
		{
			name: "ErrClearCart",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			want: 0,
			err:  errClearCart,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(cartItems, nil)
				mock.ClearCartMock.Expect(ctx, userID).Return(errClearCart)
				return mock
			},
			orderCreatorMock: func(mc *minimock.Controller) OrderCreator {
				mock := mocks.NewOrderCreatorMock(mc)
				mock.CreateOrderMock.Expect(ctx, userID, cartItems).Return(orderID, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(cfg, tt.checkoutRepoMock(mc), tt.orderCreatorMock(mc))

			got, err := domain.CreateOrder(tt.args.ctx, tt.args.userID)
			require.Equal(t, tt.want, got)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
