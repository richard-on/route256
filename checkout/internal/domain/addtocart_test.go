package domain

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain/mocks"
	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/model"
)

func TestAddToCart(t *testing.T) {
	t.Parallel()
	type checkoutRepoMockFunc func(mc *minimock.Controller) CheckoutRepo
	type stockCheckerMockFunc func(mc *minimock.Controller) StockChecker

	type args struct {
		ctx  context.Context
		user int64
		item model.Item
	}

	var (
		mc     = minimock.NewController(t)
		ctx    = context.Background()
		userID = gofakeit.Int64()
		sku    = gofakeit.Uint32()
		count  = gofakeit.Number(1, 5)

		item = model.Item{
			SKU:   sku,
			Count: uint16(count),
		}

		smallStocks = []*model.Stock{
			{
				WarehouseID: 1,
				Count:       1,
			},
		}

		stockErr     = errors.New("stock error")
		cartCountErr = errors.New("cart count error")
		addToCartErr = errors.New("add to cart error")

		stocksLen = 5
		stocks    []*model.Stock
	)

	t.Cleanup(mc.Finish)

	for i := 0; i < stocksLen; i++ {
		stock := &model.Stock{
			WarehouseID: int64(gofakeit.Number(1, 500)),
			Count:       uint64(gofakeit.Number(5, 50000)),
		}
		stocks = append(stocks, stock)
	}

	tests := []struct {
		name             string
		args             args
		err              error
		checkoutRepoMock checkoutRepoMockFunc
		stockCheckerMock stockCheckerMockFunc
	}{
		{
			name: "Positive-NoItemsInCart",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: nil,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctx, userID, item).Return(0, ErrNotInCart)

				mock.AddToCartMock.Expect(ctx, userID, item).Return(nil)
				return mock
			},
			stockCheckerMock: func(mc *minimock.Controller) StockChecker {
				mock := mocks.NewStockCheckerMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(stocks, nil)
				return mock
			},
		},
		{
			name: "Positive-ItemsInCart",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: nil,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctx, userID, item).Return(int32(count), nil)

				mock.AddToCartMock.Expect(ctx, userID, item).Return(nil)
				return mock
			},
			stockCheckerMock: func(mc *minimock.Controller) StockChecker {
				mock := mocks.NewStockCheckerMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(stocks, nil)
				return mock
			},
		},
		{
			name: "ErrStocksCheck",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: errors.WithMessage(stockErr, "checking stock"),
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				return mocks.NewCheckoutRepoMock(mc)
			},
			stockCheckerMock: func(mc *minimock.Controller) StockChecker {
				mock := mocks.NewStockCheckerMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(nil, stockErr)
				return mock
			},
		},
		{
			name: "ErrGetItemCartCount",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: cartCountErr,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctx, userID, item).Return(0, cartCountErr)
				return mock
			},
			stockCheckerMock: func(mc *minimock.Controller) StockChecker {
				mock := mocks.NewStockCheckerMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(stocks, nil)
				return mock
			},
		},
		{
			name: "ErrInsufficientStocks",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: ErrInsufficientStocks,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctx, userID, item).
					Return(int32(stocksLen), nil)

				return mock
			},
			stockCheckerMock: func(mc *minimock.Controller) StockChecker {
				mock := mocks.NewStockCheckerMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(smallStocks, nil)
				return mock
			},
		},
		{
			name: "ErrAddToCart",
			args: args{
				ctx:  ctx,
				user: userID,
				item: item,
			},
			err: addToCartErr,
			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetItemCartCountMock.Expect(ctx, userID, item).Return(int32(count), nil)
				mock.AddToCartMock.Expect(ctx, userID, item).Return(addToCartErr)
				return mock
			},
			stockCheckerMock: func(mc *minimock.Controller) StockChecker {
				mock := mocks.NewStockCheckerMock(mc)
				mock.StocksMock.Expect(ctx, sku).Return(stocks, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(tt.checkoutRepoMock(mc), tt.stockCheckerMock(mc))

			err := domain.AddToCart(tt.args.ctx, tt.args.user, tt.args.item)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
