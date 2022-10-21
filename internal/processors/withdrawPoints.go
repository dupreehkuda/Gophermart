package processors

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strconv"

	balerr "github.com/dupreehkuda/Gophermart/internal"
)

func (p processors) WithdrawPoints(order string, sum decimal.Decimal) error {
	orderConv, err := strconv.Atoi(order)
	if err != nil {
		p.logger.Error("Error occurred converting string order to int", zap.Error(err))
		return err
	}

	valid := luhnValid(orderConv)
	if !valid {
		p.logger.Error("Order number is luhn invalid", zap.Error(err))
		return balerr.BalanceInvalidOrder
	}

	balanceOk, err := p.storage.CheckPoints(orderConv, sum)
	if err != nil {
		p.logger.Error("Error occurred when checking balance", zap.Error(err))
		return err
	}

	if !balanceOk {
		return balerr.BalanceNotEnoughPoints
	}

	err = p.storage.WithdrawPoints(orderConv, sum)
	if err != nil {
		p.logger.Error("Error occurred when withdrawing points", zap.Error(err))
		return err
	}

	return nil
}
