package actions

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal"
)

// Register processes user's registration interaction
func (a actions) Register(login, password string) error {
	passwordSalt, err := i.RandSymbols(10)
	if err != nil {
		a.logger.Error("Generating salt error", zap.Error(err))
		return err
	}

	passwordHash := mdHash(password, passwordSalt)

	exists, err := a.storage.CheckDuplicateUser(login)
	if err != nil {
		a.logger.Error("User check db error", zap.Error(err))
		return err
	}

	if exists {
		return i.ErrCredentialsInUse
	}

	err = a.storage.CreateUser(login, passwordHash, passwordSalt)
	if err != nil {
		a.logger.Error("User creation db error", zap.Error(err))
		return err
	}

	return nil
}

// Login processes user's login interaction
func (a actions) Login(login, password string) error {
	exists, err := a.storage.CheckDuplicateUser(login)
	if err != nil {
		a.logger.Error("User check db error", zap.Error(err))
		return err
	}

	if !exists {
		return i.ErrWrongCredentials
	}

	passwordHash, passwordSalt, err := a.storage.LoginUser(login)
	if err != nil {
		return err
	}

	checkHash := mdHash(password, passwordSalt)
	if checkHash != passwordHash {
		return nil
	}

	return nil
}

// mdHash hashes password with salt
func mdHash(password, passwordSalt string) string {
	hasher := md5.New()
	io.WriteString(hasher, password+passwordSalt)

	return hex.EncodeToString(hasher.Sum(nil))
}
