package actions

func (a actions) GetBalance(login string) ([]byte, error) {
	data, err := a.storage.GetBalance(login)
	if err != nil {
		return nil, err
	}

	return data, nil
}
