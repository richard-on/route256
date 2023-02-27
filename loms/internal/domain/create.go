package domain

import "context"

type Item struct {
	SKU   uint32
	Count uint16
}

func (d *Domain) CreateOrder(ctx context.Context, user int64, items []Item) (int64, error) {

	// Blank business logic

	return 42, nil
}
