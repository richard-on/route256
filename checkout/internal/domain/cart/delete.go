package cart

import "context"

type Deleter interface {
	Delete(ctx context.Context, user int64, sku uint32, count uint16)
}

type Delete struct {
	deleter Deleter
}

func NewDeleter() *Delete {
	return &Delete{}
}
