package actions

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strconv"

	balanceError "github.com/dupreehkuda/Gophermart/internal"
)

func (a actions) WithdrawPoints(login, order string, sum decimal.Decimal) error {
	orderConv, err := strconv.Atoi(order)
	if err != nil {
		a.logger.Error("Error occurred converting string order to int", zap.Error(err))
		return err
	}

	valid := luhnValid(orderConv)
	if !valid {
		a.logger.Error("Order number is luhn invalid", zap.Error(err))
		return balanceError.BalanceInvalidOrderError
	}

	balanceOk, err := a.storage.CheckPoints(login, sum)
	if err != nil {
		a.logger.Error("Error occurred when checking balance", zap.Error(err))
		return err
	}

	if !balanceOk {
		return balanceError.BalanceNotEnoughPointsError
	}

	err = a.storage.WithdrawPoints(login, orderConv, sum)
	if err != nil {
		a.logger.Error("Error occurred when withdrawing points", zap.Error(err))
		return err
	}

	return nil
}
