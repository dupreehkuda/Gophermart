package sqlxpq

import (
	"github.com/shopspring/decimal"
)

// UpdateAccrual updates order data in the database
func (s storageLpq) UpdateAccrual(order int, status string, accrual decimal.Decimal) error {
	tx := s.conn.MustBegin()

	tx.MustExec("update orders set status = $1, accrual = accrual + $2 where orderid = $3;", status, accrual, order)
	tx.MustExec("update accrual set points = points + $1 where login = (select login from orders where orderid = $2);", accrual, order)

	err := tx.Commit()
	if err != nil {
		s.logger.Error("Error occurred while updating accrual data")
		return err
	}

	return nil
}
