package storage

import (
	"context"

	"github.com/jackc/pgx/v5"

	i "github.com/dupreehkuda/Gophermart/internal"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CheckPoints checks if user have enough points
func (s storage) CheckPoints(login string, sum decimal.Decimal) (decimal.Decimal, error) {
	var currentPoints decimal.Decimal

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return currentPoints, err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	err = conn.QueryRow(context.Background(), "select points from accrual where login = $1;", login).Scan(&currentPoints)
	if err != nil {
		s.logger.Error("Error occurred while scanning", zap.Error(err))
	}

	if currentPoints.LessThan(sum) {
		return currentPoints, i.ErrBalanceNotEnoughPoints
	}

	return currentPoints, nil
}

// WithdrawPoints withdraws points from users account and sets order as paid
func (s storage) WithdrawPoints(login string, order int, sum, current decimal.Decimal) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	batch := &pgx.Batch{}

	batch.Queue("update orders set pointsspent = $1 where orderid = $2;", true, order)
	batch.Queue("update accrual set points = $1 where login = $2;", current.Sub(sum), login)
	batch.Queue("update accrual set withdrawn = withdrawn + $1 where login = $2;", sum, login)

	br := conn.SendBatch(context.Background(), batch)
	defer s.batchClosing(br)

	return nil
}
