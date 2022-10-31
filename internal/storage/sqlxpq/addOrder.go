package sqlxpq

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

// CheckOrderExistence checks if order already exist in the database
func (s storageLpq) CheckOrderExistence(login string, orderID int) (bool, bool) {
	var dbUserLogin string

	s.conn.QueryRow("select login from orders where orderid = $1;", orderID).Scan(&dbUserLogin)

	if dbUserLogin == "" {
		return false, true
	}

	if dbUserLogin != login {
		return true, false
	}

	return true, true
}

// NewOrder creates inserts new order in the database
func (s storageLpq) NewOrder(login string, orderID int) error {
	rows, err := s.conn.Query("insert into orders(login, orderid, orderdate, status, accrual, pointsspent) values ($1, $2, $3, 'NEW', $4, $5);", login, orderID, time.Now().Format("2006-01-02 15:04:05"), decimal.Zero, false)
	if err != nil || rows.Err() != nil {
		s.logger.Error("Error while inserting order", zap.Error(err))
		return err
	}

	return nil
}
