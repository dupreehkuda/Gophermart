package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Register handles an action of creating new user
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

	err = h.actions.Register(regCredit.Login, regCredit.Password)

	switch {
	case errors.Is(err, i.ErrCredentialsInUse):
		h.logger.Info("Login occupied", zap.String("login", regCredit.Login))
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		h.logger.Error("Unable to call actions", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := i.GenerateJWT(regCredit.Login)
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
