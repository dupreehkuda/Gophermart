package handlers

import (
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

func (h handlers) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	var ctxKey i.LoginKey = "login"
	login := r.Context().Value(ctxKey).(string)

	response, err := h.actions.GetWithdrawals(login)
	if err != nil {
		h.logger.Error("Error call when getting withdrawals", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		h.logger.Error("Error occurred writing response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
