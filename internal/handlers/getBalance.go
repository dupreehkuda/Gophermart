package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

func (h handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login").(string)

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
