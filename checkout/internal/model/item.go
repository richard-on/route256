// Package model defines main models for business logic.
package model

// ProductInfo represents product's name and price.
type ProductInfo struct {
	// Name of the product.
	Name string
	// Price of a single product.
	Price uint32
}

// Item represents a product to buy.
type Item struct {
	// SKU is the product's stock keeping unit.
	SKU uint32
	// Count is the number of product's with this SKU.
	Count uint16
	// ProductInfo stores additional info about Item.
	ProductInfo ProductInfo
}
