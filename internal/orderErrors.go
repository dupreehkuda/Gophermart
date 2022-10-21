package internal

import (
	"errors"
)

var (
	OrderInvalidNum = errors.New("invalid order number")
	OrderOccupied   = errors.New("occupied order")
	OrderUploaded   = errors.New("already uploaded")
)

//type AddOrderError struct {
//	Err error
//}
//
//func (r *AddOrderError) Error() string {
//	return fmt.Sprintf("err %v", r.Err)
//}
