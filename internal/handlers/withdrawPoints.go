package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

type withdrawData struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}

func (h handlers) WithdrawPoints(w http.ResponseWriter, r *http.Request) {
	var ctxKey i.LoginKey = "login"
	login := r.Context().Value(ctxKey).(string)

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

	err = h.actions.WithdrawPoints(login, data.Order, data.Sum)

	switch err {
	case i.BalanceNotEnoughPointsError:
		w.WriteHeader(http.StatusPaymentRequired)
		return
	case i.OrderInvalidNumError:
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	case nil:
		return
	default:
		h.logger.Error("Error call to actions", zap.Error(err))
		return
	}
}
