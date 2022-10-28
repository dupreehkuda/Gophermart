package handlers

import (
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

func (h handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	var ctxKey i.LoginKey = "login"
	login := r.Context().Value(ctxKey).(string)

	data, err := h.actions.GetOrders(login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("Getting order error", zap.Error(err))
		return
	}

	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		h.logger.Debug("No data to return")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("Writing body error", zap.Error(err))
		return
	}
}
