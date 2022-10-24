package storage

import (
	"context"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s storage) UpdateAccrual(order int, status string, accrual decimal.Decimal) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	batch := &pgx.Batch{}

	batch.Queue("update orders set status = $1, accrual = accrual + $2 where orderid = $3;", status, accrual, order)
	batch.Queue("update accrual set points = points + $1 where login = (select login from orders where orderid = $2);", accrual, order)

	br := conn.SendBatch(context.Background(), batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			s.logger.Error("Error while updating accrual", zap.Error(err))
			return
		}
	}(br)

	s.logger.Debug("updating accrual", zap.String("status", status), zap.Float64("accrual", accrual.InexactFloat64()))

	return nil
}
