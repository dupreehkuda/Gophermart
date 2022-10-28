package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

// Login handles an action of user logging in
func (h handlers) Login(w http.ResponseWriter, r *http.Request) {
	var logCredit Credentials

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&logCredit)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if logCredit.Login == "" && logCredit.Password == "" {
		h.logger.Info("Credentials empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.actions.Login(logCredit.Login, logCredit.Password)
	if err != nil {
		h.logger.Error("Unable to authorize", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if errors.Is(err, i.ErrWrongCredentials) {
		h.logger.Error("Login or password is wrong", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := i.GenerateJWT(logCredit.Login)
	if err != nil {
		h.logger.Error("Error while generating jwt", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "JWT",
		Value: token,
	})
}
