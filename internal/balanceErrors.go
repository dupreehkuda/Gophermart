package internal

import "errors"

var (
	ErrBalanceInvalidOrder    = errors.New("invalid order number")
	ErrBalanceNotEnoughPoints = errors.New("not enough points")
)
