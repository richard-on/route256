// Package schema defines main models for repository layer..
package schema

// Item represents a product as it is stored in database.
type Item struct {
	// SKU is the product's stock keeping unit.
	SKU int64 `db:"sku"`
	// Count is the number of product's with this SKU.
	Count int32 `db:"count"`
}
