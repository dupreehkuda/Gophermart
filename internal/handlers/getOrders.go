package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

func (h handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login").(string)

	data, err := h.processor.GetOrders(login)
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("Writing body error", zap.Error(err))
		return
	}
}
