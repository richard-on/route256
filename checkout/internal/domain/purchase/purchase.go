package purchase

import "context"

type Orderer interface {
	Order(ctx context.Context, user int64) error
}

type Order struct {
	orderer Orderer
}

func New(orderer Orderer) *Order {
	return &Order{
		orderer: orderer,
	}
}

func (o *Order) Create(ctx context.Context, user int64) error {
	err := o.orderer.Order(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
