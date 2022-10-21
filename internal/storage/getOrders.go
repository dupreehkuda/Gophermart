package storage

import (
	"context"
	"encoding/json"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"time"

	"github.com/jackc/pgtype"
	"go.uber.org/zap"
)

type dbOrder struct {
	Number     string           `db:"orderID"`
	Status     string           `db:"status"`
	Accrual    decimal.Decimal  `db:"accrual"`
	UploadedAt pgtype.Timestamp `db:"orderdate"`
}

type order struct {
	Number     string          `json:"number"`
	Status     string          `json:"status"`
	Accrual    decimal.Decimal `json:"accrual,omitempty"`
	UploadedAt string          `json:"uploaded_at"`
}

func (s storage) GetOrders(login string) ([]byte, error) {
	var dataFromDB []dbOrder
	var data []order

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	rows, err := conn.Query(context.Background(), "select orderid, status, accrual, orderdate from orders where login = $1 order by orderdate;", login)
	if err != nil {
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

	// todo: figure out how to write without empty fields
	for _, val := range dataFromDB {
		data = append(data, order{
			Number:     val.Number,
			Status:     val.Status,
			Accrual:    val.Accrual,
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
