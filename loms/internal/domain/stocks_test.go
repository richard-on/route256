package domain

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/domain/mocks"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"testing"
)

func TestStocks(t *testing.T) {
	t.Parallel()
	type LOMSRepoMockFunc func(mc *minimock.Controller) LOMSRepo

	type args struct {
		ctx context.Context
		sku uint32
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()
		sku = gofakeit.Uint32()

		stocks = []model.Stock{
			{
				WarehouseID: int64(gofakeit.Number(1, 100)),
				Count:       uint64(gofakeit.Number(1, 100000)),
			},
		}

		stockErr = errors.New("stock error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         []model.Stock
		err          error
		lomsRepoMock LOMSRepoMockFunc
	}{
		{
			name: "Positive",
			args: args{
				ctx: ctx,
				sku: sku,
			},
			want: stocks,
			err:  nil,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.GetStocksMock.Expect(ctx, sku).Return(stocks, nil)
				return mock
			},
		},
		{
			name: "ErrStocksRepo",
			args: args{
				ctx: ctx,
				sku: sku,
			},
			want: nil,
			err:  stockErr,
			lomsRepoMock: func(mc *minimock.Controller) LOMSRepo {
				mock := mocks.NewLOMSRepoMock(mc)
				mock.GetStocksMock.Expect(ctx, sku).Return(nil, stockErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			domain := NewMockDomain(tt.lomsRepoMock(mc))

			res, err := domain.Stocks(tt.args.ctx, tt.args.sku)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
