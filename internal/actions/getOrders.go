package actions

// GetOrders gets a list of users completed orders
func (a actions) GetOrders(login string) ([]byte, error) {
	data, err := a.storage.GetOrders(login)
	if err != nil {
		return nil, err
	}

	return data, nil
}
