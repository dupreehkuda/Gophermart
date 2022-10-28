package storage

import (
	"context"

	"github.com/jackc/pgx/v5"

	i "github.com/dupreehkuda/Gophermart/internal"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s storage) CheckPoints(login string, sum decimal.Decimal) error {
	var currentPoints decimal.Decimal

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	err = conn.QueryRow(context.Background(), "select points from accrual where login = $1;", login).Scan(&currentPoints)
	if err != nil {
		s.logger.Error("Error occurred while scanning", zap.Error(err))
	}

	if currentPoints.LessThan(sum) {
		return i.BalanceNotEnoughPointsError
	}

	return nil
}

func (s storage) WithdrawPoints(login string, order int, sum decimal.Decimal) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	batch := &pgx.Batch{}

	batch.Queue("update accrual set points = points - $1, withdrawn = withdrawn + $1 where login = $2;", sum, login)
	batch.Queue("update orders set pointsspent = true where orderid = $1;", order)

	br := conn.SendBatch(context.Background(), batch)
	defer s.batchClosing(br)

	return nil
}
