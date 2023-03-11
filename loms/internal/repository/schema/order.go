package schema

import "time"

// Order represents information about the order, as it is stored in database.
type Order struct {
	// OrderID is the unique ID of this order.
	OrderID int64 `db:"order_id"`
	// UserID is the unique ID of the user, who made this order.
	UserID int64 `db:"user_id"`
	// Items is a slice of all Item in this order.
	Items []Item `db:"items"`
	// Status of order payment.
	Status int16 `db:"status"`
	// CreatedAt is the time this order was created.
	CreatedAt time.Time `db:"created_at"`
}
