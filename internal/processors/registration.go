package processors

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	"go.uber.org/zap"

	"github.com/dupreehkuda/Gophermart/internal"
)

func (p processors) Register(login, password string) (string, bool, error) {
	passwordSalt, err := internal.RandSymbols(10)
	if err != nil {
		p.logger.Error("Generating salt error", zap.Error(err))
		return "", false, err
	}

	passwordHash := mdHash(password, passwordSalt)

	exists, err := p.storage.CheckUser(login)
	if err != nil {
		p.logger.Error("User check db error", zap.Error(err))
		return "", false, err
	}

	if exists {
		return "", true, nil
	}

	err = p.storage.CreateUser(login, passwordHash, passwordSalt)
	if err != nil {
		p.logger.Error("User creation db error", zap.Error(err))
		return "", false, err
	}

	return passwordSalt, false, nil
}

func (p processors) Login(login, password string) (string, bool, error) {
	exists, err := p.storage.CheckUser(login)
	if err != nil {
		p.logger.Error("User check db error", zap.Error(err))
		return "", false, err
	}

	if !exists {
		p.logger.Error("User do not exist", zap.Error(err))
		return "", false, nil
	}

	passwordHash, passwordSalt, err := p.storage.LoginUser(login)
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
