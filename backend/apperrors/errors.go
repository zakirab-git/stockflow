package apperrors

import "errors"

var (
	ErrInsufficientStock = errors.New("INSUFFICIENT_STOCK")
	ErrUnknownSKU        = errors.New("UNKNOWN_SKU")
	ErrDuplicateSKU      = errors.New("DUPLICATE_SKU")
)