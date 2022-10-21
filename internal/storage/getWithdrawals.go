package storage

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgtype"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

type dbRespWithdrawal struct {
	Order       string
	Sum         decimal.Decimal
	ProcessedAt pgtype.Timestamp
}

type respWithdrawal struct {
	Order       string          `json:"order"`
	Sum         decimal.Decimal `json:"sum"`
	ProcessedAt string          `json:"processed_at"`
}

func (s storage) GetWithdrawals(login string) ([]byte, error) {
	dbResp := []dbRespWithdrawal{}
	resp := []respWithdrawal{}

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	rows, err := conn.Query(context.Background(), "select orderid, accrual, orderdate from orders where login = $1 and pointsspent = true order by orderdate", login)
	if err != nil {
		s.logger.Error("Error while getting withdrawals", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r dbRespWithdrawal
		err := rows.Scan(&r.Order, &r.Sum, &r.ProcessedAt)
		if err != nil {
			s.logger.Error("Error while scanning rows", zap.Error(err))
			return nil, err
		}
		dbResp = append(dbResp, r)
	}

	for _, val := range dbResp {
		resp = append(resp, respWithdrawal{
			Order:       val.Order,
			Sum:         val.Sum,
			ProcessedAt: val.ProcessedAt.Time.Format(time.RFC3339),
		})
	}

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error marshaling data", zap.Error(err))
		return nil, err
	}

	return resultJSON, nil
}
