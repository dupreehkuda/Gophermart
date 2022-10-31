package sqlxpq

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

// CheckPoints checks if user have enough points
func (s storageLpq) CheckPoints(login string, sum decimal.Decimal) (decimal.Decimal, error) {
	var currentPoints decimal.Decimal

	err := s.conn.QueryRow("select points from accrual where login = $1;", login).Scan(&currentPoints)
	if err != nil {
		s.logger.Error("Error occurred while scanning", zap.Error(err))
	}

	if currentPoints.LessThan(sum) {
		return currentPoints, i.ErrBalanceNotEnoughPoints
	}

	return currentPoints, nil
}

// WithdrawPoints withdraws points from users account and sets order as paid
func (s storageLpq) WithdrawPoints(login string, order int, sum, current decimal.Decimal) error {
	tx := s.conn.MustBegin()

	tx.MustExec("update orders set pointsspent = $1 where orderid = $2;", true, order)
	tx.MustExec("update accrual set points = points - $1, withdrawn = withdrawn + $1 where login = $2;", sum, login)

	err := tx.Commit()
	if err != nil {
		s.logger.Error("Error occurred while withdrawing data")
		return err
	}

	return nil
}
