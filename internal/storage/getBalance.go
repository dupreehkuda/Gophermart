package storage

import (
	"context"
	"encoding/json"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type respBalance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}

func (s storage) GetBalance(login string) ([]byte, error) {
	var resp respBalance

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	conn.QueryRow(context.Background(), "select points, withdrawn from accrual where login = $1;", login).Scan(&resp.Current, &resp.Withdrawn)

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error marshaling data", zap.Error(err))
		return nil, err
	}

	return resultJSON, nil
}
