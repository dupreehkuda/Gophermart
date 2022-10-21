package processors

import (
	orderr "github.com/dupreehkuda/Gophermart/internal"
)

func (p processors) NewOrder(login string, order int) error {
	valid := luhnValid(order)
	if !valid {
		return orderr.OrderInvalidNum
	}

	exists, usersOrder, err := p.storage.CheckOrder(login, order)
	if err != nil {
		return err
	}

	if exists {
		if !usersOrder {
			return orderr.OrderOccupied
		}
		return orderr.OrderUploaded
	}

	status, accrual, err := p.getAccrualData(order)
	if err != nil {
		return err
	}

	err = p.storage.NewOrder(login, status, order, accrual)
	if err != nil {
		return err
	}

	return nil
}
