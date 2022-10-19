package storage

import (
	"context"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s storage) CheckOrder(login string, order int) (bool, bool, error) {
	var dbUserLogin string

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false, false, err
	}
	defer conn.Release()

	conn.QueryRow(context.Background(), "select login from orders where orderid = $1;", order).Scan(&dbUserLogin)

	if dbUserLogin == "" {
		return false, true, nil
	}

	if dbUserLogin != login {
		return true, false, nil
	}

	return true, true, nil
}

func (s storage) NewOrder(login, status string, order int, accrual decimal.Decimal) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	_, err = conn.Query(context.Background(), "insert into orders(login, orderid, orderdate, status, accrual) values ($1, $2, $3, $4, $5);", login, order, time.Now().Format("2006-01-02 15:04:05"), status, accrual)
	if err != nil {
		s.logger.Error("Error while inserting order", zap.Error(err))
		return err
	}

	return nil
}
