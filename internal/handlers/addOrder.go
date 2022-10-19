package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	orderr "github.com/dupreehkuda/Gophermart/internal"

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

	err = h.processor.NewOrder(login, order)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	switch {
	case errors.Is(err, orderr.InvalidNum):
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	case errors.Is(err, orderr.Occupied):
		w.WriteHeader(http.StatusConflict)
		return
	case errors.Is(err, orderr.Uploaded):
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("Adding order error", zap.Error(err))
		return
	}
}
