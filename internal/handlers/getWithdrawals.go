package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

func (h handlers) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login").(string)

	response, err := h.actions.GetWithdrawals(login)
	if err != nil {
		h.logger.Error("Error call when getting withdrawals", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		h.logger.Error("Error occurred writing response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
