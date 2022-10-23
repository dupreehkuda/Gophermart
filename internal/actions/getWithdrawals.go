package actions

import "go.uber.org/zap"

func (a actions) GetWithdrawals(login string) ([]byte, error) {
	data, err := a.storage.GetWithdrawals(login)
	if err != nil {
		a.logger.Error("Error call when getting withdrawals", zap.Error(err))
		return nil, err
	}

	return data, nil
}
