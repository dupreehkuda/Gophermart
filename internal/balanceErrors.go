package internal

import "errors"

var (
	BalanceInvalidOrder    = errors.New("invalid order number")
	BalanceNotEnoughPoints = errors.New("not enough points")
)
