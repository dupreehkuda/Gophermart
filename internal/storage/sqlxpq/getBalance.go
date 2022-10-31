package sqlxpq

import (
	"encoding/json"
	"go.uber.org/zap"
	"math"
)

// GetBalance gets user's current balance from the database
func (s storageLpq) GetBalance(login string) ([]byte, error) {
	var resp respBalance
	var dbResp dbRespBalance

	s.conn.QueryRow("select points, withdrawn from accrual where login = $1;", login).Scan(&dbResp.Current, &dbResp.Withdrawn)

	fCurrent, _ := dbResp.Current.Float64()
	fWithdrawn, _ := dbResp.Withdrawn.Float64()
	resp = respBalance{
		Current:   math.Round(fCurrent*100) / 100,
		Withdrawn: math.Round(fWithdrawn*100) / 100,
	}

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error marshaling data", zap.Error(err))
		return nil, err
	}

	return resultJSON, nil
}
