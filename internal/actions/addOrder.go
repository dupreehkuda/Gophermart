package actions

import (
	orderError "github.com/dupreehkuda/Gophermart/internal"
)

func (a actions) NewOrder(login string, order int) error {
	valid := luhnValid(order)
	if !valid {
		return orderError.OrderInvalidNumError
	}

	exists, usersOrder, err := a.storage.CheckOrder(login, order)
	if err != nil {
		return err
	}

	if exists {
		if !usersOrder {
			return orderError.OrderOccupiedError
		}
		return orderError.OrderUploadedError
	}

	err = a.storage.NewOrder(login, order)
	if err != nil {
		return err
	}

	a.service.Channel <- order

	return nil
}
