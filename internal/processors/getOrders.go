package processors

func (p processors) GetOrders(login string) ([]byte, error) {
	data, err := p.storage.GetOrders(login)
	if err != nil {
		return nil, err
	}

	return data, nil
}
