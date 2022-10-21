package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	balerr "github.com/dupreehkuda/Gophermart/internal"
)

type withdrawData struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}

func (h handlers) WithdrawPoints(w http.ResponseWriter, r *http.Request) {
	var data withdrawData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.Order == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.processor.WithdrawPoints(data.Order, data.Sum)

	switch err {
	case balerr.BalanceNotEnoughPoints:
		w.WriteHeader(http.StatusPaymentRequired)
		return
	case balerr.OrderInvalidNum:
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	default:
		h.logger.Error("Error call to processors", zap.Error(err))
		return
	}
}
