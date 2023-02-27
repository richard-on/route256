package domain

import "context"

type Item struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

type ProductInfo struct {
	Name  string
	Price uint32
}

func (d *Domain) ListCart(ctx context.Context, user int64) ([]Item, error) {

	items := []Item{
		{
			SKU:   1076963,
			Count: 3,
		},
		{
			SKU:   1625903,
			Count: 1,
		},
	}

	for i, item := range items {
		product, err := d.productLister.GetProduct(ctx, item.SKU)
		if err != nil {
			return nil, err
		}
		items[i].Name = product.Name
		items[i].Price = product.Price
	}

	return items, nil
}
