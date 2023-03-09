package model

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
