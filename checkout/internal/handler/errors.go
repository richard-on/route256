package handler

import "github.com/pkg/errors"

var (
	ErrEmptyUser = errors.New("empty user")
	ErrEmptySKU  = errors.New("empty sku")
	ErrZeroCount = errors.New("zero count")
)
