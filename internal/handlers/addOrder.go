package handlers

import (
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func (h handlers) AddOrder(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login").(string)

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error("Reading body error", zap.Error(err))
		return
	}

	order, err := strconv.Atoi(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("Conversion error", zap.Error(err))
		return
	}

	valid, exists, usersOrder, err := h.processor.NewOrder(login, order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("Adding order error", zap.Error(err))
		return
	}

	switch {
	case !valid:
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	case exists && !usersOrder:
		w.WriteHeader(http.StatusConflict)
		return
	case exists && usersOrder:
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}
