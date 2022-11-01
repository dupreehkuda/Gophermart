package sqlxpq

import (
	"encoding/json"
	"math"
	"time"

	"go.uber.org/zap"
)

// GetOrders gets user's completed orders from the database
func (s storageLpq) GetOrders(login string) ([]byte, error) {
	var dataFromDB []dbOrder
	var data []order

	rows, err := s.conn.Query("select orderid, status, accrual, orderdate from orders where login = $1 order by orderdate;", login)
	if err != nil || rows.Err() != nil {
		s.logger.Error("Error while getting orders", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d dbOrder
		err := rows.Scan(&d.Number, &d.Status, &d.Accrual, &d.UploadedAt)
		if err != nil {
			s.logger.Error("Error while scanning rows", zap.Error(err))
			return nil, err
		}
		dataFromDB = append(dataFromDB, d)
	}

	for _, val := range dataFromDB {
		f, _ := val.Accrual.Float64()
		data = append(data, order{
			Number:     val.Number,
			Status:     val.Status,
			Accrual:    math.Round(f*100) / 100,
			UploadedAt: val.UploadedAt.Time.Format(time.RFC3339),
		})
	}

	if len(data) == 0 {
		return nil, nil
	}

	resultJSON, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("Error marshaling data", zap.Error(err))
		return nil, err
	}

	return resultJSON, nil
}
