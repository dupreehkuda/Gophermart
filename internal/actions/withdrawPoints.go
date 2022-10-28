package actions

import (
	"errors"
	"strconv"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

func (a actions) WithdrawPoints(login, orderID string, sum decimal.Decimal) error {
	orderConv, err := strconv.Atoi(orderID)
	if err != nil {
		a.logger.Error("Error occurred converting string order to int", zap.Error(err))
		return err
	}

	valid := luhnValid(orderConv)
	if !valid {
		a.logger.Error("Order number is luhn invalid", zap.Error(err))
		return i.ErrBalanceInvalidOrder
	}

	err = a.storage.CheckPoints(login, sum)

	switch {
	case errors.Is(err, i.ErrBalanceNotEnoughPoints):
		return err
	case err != nil:
		a.logger.Error("Error occurred when checking balance", zap.Error(err))
		return err
	}

	err = a.storage.WithdrawPoints(login, orderConv, sum)
	if err != nil {
		a.logger.Error("Error occurred when withdrawing points", zap.Error(err))
		return err
	}

	return nil
}
