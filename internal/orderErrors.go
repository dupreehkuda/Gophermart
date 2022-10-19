package internal

import (
	"errors"
)

var (
	InvalidNum = errors.New("invalid order number")
	Occupied   = errors.New("occupied order")
	Uploaded   = errors.New("already uploaded")
)

//type AddOrderError struct {
//	Err error
//}
//
//func (r *AddOrderError) Error() string {
//	return fmt.Sprintf("err %v", r.Err)
//}
