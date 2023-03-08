package loms

import (
	"math"

	"github.com/pkg/errors"
)

var (
	// ErrEmptyUser is the error returned when user ID is zero.
	ErrEmptyUser = errors.New("empty or zero user")
	// ErrEmptySKU is the error returned when sku is zero.
	ErrEmptySKU = errors.New("empty or zero sku")
	// ErrZeroCount is the error returned when count is zero.
	ErrZeroCount = errors.New("zero count")
	// ErrEmptyOrder is the error returned when order ID is zero.
	ErrEmptyOrder = errors.New("empty or zero orderID")
	// ErrCountOverflow is the error returned when value is more than max value for uint16.
	ErrCountOverflow = errors.New("uint16 overflow")
)

// validateCount validates value of count.
//
// Count should be non-zero and not bigger than 65535 (1<<16 - 1).
func validateCount(count uint32) error {
	if count == 0 {
		return ErrZeroCount
	}
	if count > math.MaxUint16 {
		return ErrCountOverflow
	}

	return nil
}

// validateUser validates value of userID.
//
// UserID should be non-zero.
func validateUser(userID int64) error {
	if userID == 0 {
		return ErrEmptyUser
	}

	return nil
}

// validateSKU validates value of SKU.
//
// SKU should be non-zero.
func validateSKU(sku uint32) error {
	if sku == 0 {
		return ErrEmptySKU
	}

	return nil
}

// validateOrder validates value of orderID.
//
// orderID should be non-zero.
func validateOrder(orderID int64) error {
	if orderID == 0 {
		return ErrEmptyOrder
	}

	return nil
}
