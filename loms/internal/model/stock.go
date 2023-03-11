package model

// Stock represents a number of specific product available in a specific warehouse.
type Stock struct {
	// WarehouseID is the ID of a warehouse where the item is stored.
	WarehouseID int64
	// Count is the number of specific product available in this warehouse.
	Count uint64
}
