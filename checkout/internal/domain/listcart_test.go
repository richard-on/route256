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

func TestListCart(t *testing.T) {
	t.Parallel()
	type checkoutRepoMockFunc func(mc *minimock.Controller) CheckoutRepo
	type productListerMockFunc func(mc *minimock.Controller) ProductLister

	type args struct {
		ctx  context.Context
		user int64
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		userID int64 = 1

		itemNum    = 5
		cartItems  []model.Item
		products   []model.ProductInfo
		totalPrice uint32

		errListProduct = errors.New("error while listing products")
	)

	t.Cleanup(mc.Finish)

	for i := 0; i < itemNum; i++ {
		name := gofakeit.CarModel()
		sku := uint32(gofakeit.Number(50000, 5000000))
		count := uint32(gofakeit.Number(1, 100))
		price := uint32(gofakeit.Number(100, 10000))
		totalPrice += price * count

		cartItems = append(cartItems, model.Item{
			SKU:   sku,
			Count: uint16(count),
		})

		products = append(products, model.ProductInfo{
			Name:  name,
			Price: price,
		})
	}

	tests := []struct {
		name              string
		args              args
		sku               uint32
		wantItems         []model.Item
		wantTotalPrice    uint32
		err               error
		checkoutRepoMock  checkoutRepoMockFunc
		productListerMock productListerMockFunc
	}{
		{
			name: "Positive",
			args: args{
				ctx:  ctx,
				user: userID,
			},
			wantItems:      cartItems,
			wantTotalPrice: totalPrice,
			err:            nil,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(cartItems, nil)
				return mock
			},
			productListerMock: func(mc *minimock.Controller) ProductLister {
				mock := mocks.NewProductListerMock(mc)
				for i := 0; i < itemNum; i++ {
					mock.GetProductMock.When(context.Background(), cartItems[i].SKU).Then(products[i], nil)
				}

				return mock
			},
		},
		{
			name: "ErrEmptyCart",
			args: args{
				ctx:  ctx,
				user: userID,
			},
			wantItems:      nil,
			wantTotalPrice: 0,
			err:            ErrEmptyCart,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(nil, ErrEmptyCart)
				return mock
			},
			productListerMock: func(mc *minimock.Controller) ProductLister {
				return mocks.NewProductListerMock(mc)
			},
		},
		{
			name: "ErrProductList",
			args: args{
				ctx:  ctx,
				user: userID,
			},
			wantItems:      nil,
			wantTotalPrice: 0,
			err:            errListProduct,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(cartItems, nil)
				return mock
			},
			productListerMock: func(mc *minimock.Controller) ProductLister {
				mock := mocks.NewProductListerMock(mc)
				mock.GetProductMock.When(context.Background(), cartItems[0].SKU).
					Then(model.ProductInfo{}, errListProduct)

				return mock
			},
		},
		{
			name: "ErrProductList-LastProduct",
			args: args{
				ctx:  ctx,
				user: userID,
			},
			wantItems:      nil,
			wantTotalPrice: 0,
			err:            errListProduct,

			checkoutRepoMock: func(mc *minimock.Controller) CheckoutRepo {
				mock := mocks.NewCheckoutRepoMock(mc)
				mock.GetCartItemsMock.Expect(ctx, userID).Return(cartItems, nil)
				return mock
			},
			productListerMock: func(mc *minimock.Controller) ProductLister {
				mock := mocks.NewProductListerMock(mc)
				for i := 0; i < itemNum-1; i++ {
					mock.GetProductMock.When(context.Background(), cartItems[i].SKU).Then(products[i], nil)
				}
				mock.GetProductMock.When(context.Background(), cartItems[itemNum-1].SKU).
					Then(model.ProductInfo{}, errListProduct)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(tt.checkoutRepoMock(mc), tt.productListerMock(mc))

			items, total, err := domain.ListCart(tt.args.ctx, tt.args.user)
			require.Equal(t, tt.wantItems, items)
			require.Equal(t, tt.wantTotalPrice, total)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}

}
