package sqlxpq

import (
	"encoding/json"
	"math"
	"time"

	"go.uber.org/zap"
)

// GetWithdrawals gets user's completed withdrawals from the database
func (s storageLpq) GetWithdrawals(login string) ([]byte, error) {
	dbResp := []dbWithdrawal{}
	resp := []withdrawal{}

	rows, err := s.conn.Query("select orderid, withdrawn, processed_at from withdrawals where login = $1 order by processed_at;", login)
	if err != nil || rows.Err() != nil {
		s.logger.Error("Error while getting withdrawals", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r dbWithdrawal
		err := rows.Scan(&r.Order, &r.Sum, &r.ProcessedAt)
		if err != nil {
			s.logger.Error("Error while scanning rows", zap.Error(err))
			return nil, err
		}
		dbResp = append(dbResp, r)
	}

	for _, val := range dbResp {
		f, _ := val.Sum.Float64()
		resp = append(resp, withdrawal{
			Order:       val.Order,
			Sum:         math.Round(f*100) / 100,
			ProcessedAt: val.ProcessedAt.Time.Format(time.RFC3339),
		})
	}

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error marshaling data", zap.Error(err))
		return nil, err
	}

	return resultJSON, nil
}
