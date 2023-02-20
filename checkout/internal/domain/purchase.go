package domain

import "context"

func (d *Domain) CreateOrder(ctx context.Context, user int64) error {
	err := d.orderCreator.Order(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
