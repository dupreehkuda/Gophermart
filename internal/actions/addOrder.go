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

	status, accrual, err := a.getAccrualData(order)
	if err != nil {
		return err
	}

	err = a.storage.NewOrder(login, status, order, accrual)
	if err != nil {
		return err
	}

	return nil
}
