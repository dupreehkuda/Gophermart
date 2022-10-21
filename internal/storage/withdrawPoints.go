package storage

import (
	"context"
	"github.com/jackc/pgx/v5"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s storage) CheckPoints(order int, sum decimal.Decimal) (bool, error) {
	var currentPoints decimal.Decimal

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false, err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	conn.QueryRow(context.Background(), "select points from accrual where login = (select login from orders where orderid = $1);", order).Scan(&currentPoints)
	if err != nil {
		s.logger.Error("Error while checking points", zap.Error(err))
		return false, err
	}

	if currentPoints.LessThan(sum) {
		return false, nil
	}

	return true, nil
}

func (s storage) WithdrawPoints(order int, sum decimal.Decimal) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	batch := &pgx.Batch{}

	batch.Queue("update accrual set points = points - $1, withdrawn = withdrawn + $1 where login = (select login from orders where orderid = $2);", sum, order)
	batch.Queue("update orders set pointsspent = true where orderid = $1", order)

	br := conn.SendBatch(context.Background(), batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			s.logger.Error("Error while updating", zap.Error(err))
			return
		}
	}(br)

	return nil
}