package handlers

import (
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

// GetBalance handles action of getting current balance
func (h handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	var ctxKey i.LoginKey = "login"
	login := r.Context().Value(ctxKey).(string)

	response, err := h.actions.GetBalance(login)
	if err != nil {
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
