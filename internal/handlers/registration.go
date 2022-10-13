package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	"github.com/dupreehkuda/Gophermart/internal"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h handlers) Register(w http.ResponseWriter, r *http.Request) {
	var regCredit Credentials

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&regCredit)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if regCredit.Login == "" && regCredit.Password == "" {
		h.logger.Info("Credentials empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, occupied, err := h.processor.Register(regCredit.Login, regCredit.Password)
	if err != nil {
		h.logger.Error("Unable to call processor", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if occupied {
		h.logger.Info("Login occupied", zap.String("login", regCredit.Login))
		w.WriteHeader(http.StatusConflict)
		return
	}

	token, err := internal.GenerateJWT(regCredit.Login)
	if err != nil {
		h.logger.Error("Error while generating jwt", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "JWT",
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}
