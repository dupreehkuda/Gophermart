package processors

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
		p.logger.Debug("If exists return in processors")
		return valid, exists, usersOrder, nil
	}
	
	err = p.storage.NewOrder(login, order)
	if err != nil {
		return valid, false, false, err
	}

	return valid, exists, usersOrder, nil
}
