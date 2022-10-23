package storage

import (
	"context"
	"encoding/json"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"go.uber.org/zap"
)

func (s storage) GetBalance(login string) ([]byte, error) {
	var resp respBalance
	var dbResp dbRespBalance

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	conn.QueryRow(context.Background(), "select points, withdrawn from accrual where login = $1;", login).Scan(&dbResp.Current, &dbResp.Withdrawn)

	fCurrent, _ := dbResp.Current.Float64()
	fWithdrawn, _ := dbResp.Withdrawn.Float64()
	resp = respBalance{
		Current:   float32(fCurrent),
		Withdrawn: float32(fWithdrawn),
	}

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error marshaling data", zap.Error(err))
		return nil, err
	}

	return resultJSON, nil
}
