package domain

import "context"

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Info struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
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
