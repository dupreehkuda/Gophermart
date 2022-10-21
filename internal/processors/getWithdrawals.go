package processors

import "go.uber.org/zap"

func (p processors) GetWithdrawals(login string) ([]byte, error) {
	data, err := p.storage.GetWithdrawals(login)
	if err != nil {
		p.logger.Error("Error call when getting withdrawals", zap.Error(err))
		return nil, err
	}

	return data, nil
}
