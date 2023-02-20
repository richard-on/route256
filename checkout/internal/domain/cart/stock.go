package cart

import (
	"context"
)

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Checker interface {
	Stock(ctx context.Context, sku uint32) ([]Stock, error)
}

type Check struct {
	checker Checker
}

func NewChecker(checker Checker) *Check {
	return &Check{
		checker: checker,
	}
}
