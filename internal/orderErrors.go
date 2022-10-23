package internal

import (
	"errors"
)

var (
	OrderInvalidNumError = errors.New("invalid order number")
	OrderOccupiedError   = errors.New("occupied order")
	OrderUploadedError   = errors.New("already uploaded")
)
