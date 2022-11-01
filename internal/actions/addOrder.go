package actions

import (
	orderError "github.com/dupreehkuda/Gophermart/internal"
)

// NewOrder is processing new orders in the system
func (a actions) NewOrder(login string, orderID int) error {
	valid := luhnValid(orderID)
	if !valid {
		return orderError.ErrOrderInvalidNum
	}

	exists, usersOrder := a.storage.CheckOrderExistence(login, orderID)

	if exists {
		if !usersOrder {
			return orderError.ErrOrderOccupied
		}
		return orderError.ErrOrderUploaded
	}

	err := a.storage.NewOrder(login, orderID)
	if err != nil {
		return err
	}

	a.service.OrderQueue <- orderID

	return nil
}
