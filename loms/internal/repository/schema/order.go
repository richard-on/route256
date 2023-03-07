package schema

import "time"

type OrderInfo struct {
	OrderID   int64     `db:"order_id"`
	UserID    int64     `db:"user_id"`
	Items     ItemArray `db:"items"`
	Status    uint8     `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
