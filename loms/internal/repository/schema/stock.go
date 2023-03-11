package schema

// Stock represents a number of specific product available in a specific warehouse as it is stored in database.
type Stock struct {
	// WarehouseID is the ID of a warehouse where the item is stored.
	WarehouseID int64 `db:"warehouse_id"`
	// Count is the number of specific product available in this warehouse.
	Count int32 `db:"count"`
}

// Reserve is a helper struct for retrieving reserved items from a database.
type Reserve struct {
	// SKU is the product's stock keeping unit.
	SKU int64
	// WarehouseID is the ID of a warehouse where the item is reserved.
	WarehouseID int64
	// Count is the number of specific product to reserve.
	Count int64
}
