package actions

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	"go.uber.org/zap"

	"github.com/dupreehkuda/Gophermart/internal"
)

func (a actions) Register(login, password string) (string, bool, error) {
	passwordSalt, err := internal.RandSymbols(10)
	if err != nil {
		a.logger.Error("Generating salt error", zap.Error(err))
		return "", false, err
	}

	passwordHash := mdHash(password, passwordSalt)

	exists, err := a.storage.CheckUser(login)
	if err != nil {
		a.logger.Error("User check db error", zap.Error(err))
		return "", false, err
	}

	if exists {
		return "", true, nil
	}

	err = a.storage.CreateUser(login, passwordHash, passwordSalt)
	if err != nil {
		a.logger.Error("User creation db error", zap.Error(err))
		return "", false, err
	}

	return passwordSalt, false, nil
}

func (a actions) Login(login, password string) (string, bool, error) {
	exists, err := a.storage.CheckUser(login)
	if err != nil {
		a.logger.Error("User check db error", zap.Error(err))
		return "", false, err
	}

	if !exists {
		a.logger.Error("User do not exist", zap.Error(err))
		return "", false, nil
	}

	passwordHash, passwordSalt, err := a.storage.LoginUser(login)
	if err != nil {
		return "", false, err
	}

	checkHash := mdHash(password, passwordSalt)
	if checkHash != passwordHash {
		return "", false, nil
	}

	return passwordSalt, true, nil
}

func mdHash(password, passwordSalt string) string {
	hasher := md5.New()
	io.WriteString(hasher, password+passwordSalt)

	return hex.EncodeToString(hasher.Sum(nil))
}
