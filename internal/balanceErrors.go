package internal

import "errors"

var (
	BalanceInvalidOrderError    = errors.New("invalid order number")
	BalanceNotEnoughPointsError = errors.New("not enough points")
)
