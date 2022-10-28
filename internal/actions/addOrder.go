package actions

import (
	orderError "github.com/dupreehkuda/Gophermart/internal"
)

func (a actions) NewOrder(login string, orderID int) error {
	valid := luhnValid(orderID)
	if !valid {
		return orderError.OrderInvalidNumError
	}

	exists, usersOrder := a.storage.CheckOrderExistence(login, orderID)

	if exists {
		if !usersOrder {
			return orderError.OrderOccupiedError
		}
		return orderError.OrderUploadedError
	}

	err := a.storage.NewOrder(login, orderID)
	if err != nil {
		return err
	}

	a.service.OrderQueue <- orderID

	return nil
}
