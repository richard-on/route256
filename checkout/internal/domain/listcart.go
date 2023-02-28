package domain

import "context"

// Item represents a product in the cart.
type Item struct {
	// SKU is the product's stock keeping unit.
	SKU uint32
	// Count is the number of product's with this SKU in a cart.
	Count uint16
	// Name of the product.
	Name string
	// Price of a single product.
	Price uint32
}

// ProductInfo represents product's name and price.
type ProductInfo struct {
	// Name of the product.
	Name string
	// Price of the product.
	Price uint32
}

// ListCart lists all products that are currently in a user's cart.
func (d *Domain) ListCart(ctx context.Context, user int64) ([]Item, uint32, error) {
	// Example items
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

	var totalPrice uint32 = 0
	for i, item := range items {
		product, err := d.productLister.GetProduct(ctx, item.SKU)
		if err != nil {
			return nil, 0, err
		}
		items[i].Name = product.Name
		items[i].Price = product.Price
		totalPrice += product.Price
	}

	return items, totalPrice, nil
}
