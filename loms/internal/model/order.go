package model

const (
	Unspecified     Status = iota // Unspecified status.
	NewOrder                      // NewOrder is the status for a newly created order.
	AwaitingPayment               // AwaitingPayment is the status for an order that awaits payment.
	Failed                        // Failed is the status for an order whose payment has failed.
	Paid                          // Paid is the status for a successfully paid order.
	Cancelled                     // Cancelled is the status for a cancelled order.
)

// Status is an enumeration that represents a status of order payment.
type Status uint8

// Order represents information about the order, including its current Status,
// User who made the order and Items in the order.
type Order struct {
	// Status of order payment.
	Status Status
	// User is the unique ID of the user, who made this order.
	User int64
	// Items is the slice of all Item in the order.
	Items []Item
}

type OrderStatus struct {
	// ID of the order.
	ID int64 `json:"id"`
	// Status of order payment.
	Status Status `json:"status"`
}
