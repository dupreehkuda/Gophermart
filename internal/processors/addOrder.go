package processors

import "go.uber.org/zap"

func (p processors) NewOrder(login string, order int) (bool, bool, bool, error) {
	valid := luhnValid(order)
	if !valid {
		return valid, false, false, nil
	}

	exists, usersOrder, err := p.storage.CheckOrder(login, order)
	if err != nil {
		return valid, exists, usersOrder, err
	}

	if exists || !usersOrder {
		p.logger.Debug("If exists return in processors", zap.Any("Order exists?", exists), zap.Any("Users order?", usersOrder))
		return valid, exists, usersOrder, nil
	}

	status, accrual, err := p.getAccrualData(order)
	if err != nil {
		return valid, false, false, err
	}

	err = p.storage.NewOrder(login, status, order, accrual)
	if err != nil {
		return valid, false, false, err
	}

	return valid, exists, usersOrder, nil
}
