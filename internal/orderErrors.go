package internal

import (
	"errors"
)

var (
	ErrOrderInvalidNum = errors.New("invalid order number")
	ErrOrderOccupied   = errors.New("occupied order")
	ErrOrderUploaded   = errors.New("already uploaded")
)
