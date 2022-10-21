package processors

func (p processors) GetBalance(login string) ([]byte, error) {
	data, err := p.storage.GetBalance(login)
	if err != nil {
		return nil, err
	}

	return data, nil
}
