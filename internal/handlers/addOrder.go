package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	i "github.com/dupreehkuda/Gophermart/internal"

	"go.uber.org/zap"
)

// AddOrder handles action of processing new orders
func (h handlers) AddOrder(w http.ResponseWriter, r *http.Request) {
	var ctxKey i.LoginKey = "login"
	login := r.Context().Value(ctxKey).(string)

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

	err = h.actions.NewOrder(login, order)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	switch {
	case errors.Is(err, i.ErrOrderInvalidNum):
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	case errors.Is(err, i.ErrOrderOccupied):
		w.WriteHeader(http.StatusConflict)
		return
	case errors.Is(err, i.ErrOrderUploaded):
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("Adding order error", zap.Error(err))
		return
	}
}
