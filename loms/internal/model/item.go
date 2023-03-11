// Package model defines main models for business logic.
package model

// Item represents a product to buy.
type Item struct {
	// SKU is the product's stock keeping unit.
	SKU uint32
	// Count is the number of product's with this SKU.
	Count uint16
}
