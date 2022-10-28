package storage

import (
	"context"
	"encoding/json"
	"math"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"

	"go.uber.org/zap"
)

// GetWithdrawals gets user's completed withdrawals from the database
func (s storage) GetWithdrawals(login string) ([]byte, error) {
	dbResp := []dbWithdrawal{}
	resp := []withdrawal{}

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	rows, err := conn.Query(context.Background(), "select orderid, accrual, orderdate from orders where pointsspent = $1 order by orderdate;", true)
	if err != nil {
		s.logger.Error("Error while getting withdrawals", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	s.logger.Debug("after connection", zap.Any("rows", rows))

	for rows.Next() {
		var r dbWithdrawal
		err := rows.Scan(&r.Order, &r.Sum, &r.ProcessedAt)
		if err != nil {
			s.logger.Error("Error while scanning rows", zap.Error(err))
			return nil, err
		}
		s.logger.Debug("if rows cycle", zap.Any("new item", r))
		dbResp = append(dbResp, r)
	}

	for _, val := range dbResp {
		f, _ := val.Sum.Float64()
		resp = append(resp, withdrawal{
			Order:       val.Order,
			Sum:         math.Round(f*100) / 100,
			ProcessedAt: val.ProcessedAt.Time.Format(time.RFC3339),
		})
		s.logger.Debug("if redefine cycle", zap.Any("new item added", resp))
	}

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error marshaling data", zap.Error(err))
		return nil, err
	}

	return resultJSON, nil
}
