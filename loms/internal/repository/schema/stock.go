package schema

type Stock struct {
	WarehouseID int64 `db:"warehouse_id"`
	Count       int32 `db:"count"`
}
