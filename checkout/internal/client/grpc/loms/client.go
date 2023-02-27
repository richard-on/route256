package loms

import (
	"context"

	"gitlab.ozon.dev/rragusskiy/homework-1/checkout/internal/domain"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc"
)

type Client interface {
	Order(ctx context.Context, user int64) (domain.OrderInfo, error)
	Stock(ctx context.Context, sku uint32) ([]*domain.Stock, error)
}

type client struct {
	lomsClient loms.LOMSClient
}

func NewClient(cc *grpc.ClientConn) *client {
	return &client{
		lomsClient: loms.NewLOMSClient(cc),
	}
}

func (c *client) Order(ctx context.Context, user int64) (domain.OrderInfo, error) {
	items := []domain.Item{
		{
			SKU:   123456,
			Count: 1,
		},
		{
			SKU:   654321,
			Count: 2,
		},
	}

	itemsReq := make([]*loms.Item, 0, 2)
	for _, item := range items {
		itemsReq = append(itemsReq, &loms.Item{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}
	orderReq := &loms.CreateOrderRequest{
		User:  user,
		Items: itemsReq,
	}

	resp, err := c.lomsClient.CreateOrder(ctx, orderReq)
	if err != nil {
		return domain.OrderInfo{}, err
	}

	return domain.OrderInfo{
		OrderID: resp.GetOrderId(),
	}, nil
}

func (c *client) Stock(ctx context.Context, sku uint32) ([]*domain.Stock, error) {
	resp, err := c.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: sku})
	if err != nil {
		return nil, err
	}

	stocksResp := resp.GetStocks()

	stocks := make([]*domain.Stock, 0, len(stocksResp))
	for _, item := range stocksResp {
		stocks = append(stocks, &domain.Stock{
			WarehouseID: item.GetWarehouseId(),
			Count:       item.GetCount(),
		})
	}

	return stocks, nil
}
