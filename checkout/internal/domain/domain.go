package domain

import "context"

type StockChecker interface {
	Stock(ctx context.Context, sku uint32) ([]Stock, error)
}

/*type CartDeleter interface {
	Delete(ctx context.Context, user int64, sku uint32, count uint16)
}*/

type ProductLister interface {
	GetProduct(ctx context.Context, sku uint32) (Info, error)
}

type OrderCreator interface {
	Order(ctx context.Context, user int64) error
}

type Domain struct {
	stockChecker StockChecker
	//cartDeleter   CartDeleter
	productLister ProductLister
	orderCreator  OrderCreator
}

func New(checker StockChecker, lister ProductLister, creator OrderCreator) *Domain {
	return &Domain{
		stockChecker:  checker,
		orderCreator:  creator,
		productLister: lister,
	}
}
