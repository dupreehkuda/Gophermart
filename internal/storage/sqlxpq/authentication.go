package sqlxpq

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CreateUser inserts new user's data in the database
func (s storageLpq) CreateUser(login string, passwordHash string, passwordSalt string) error {
	tx := s.conn.MustBegin()

	tx.MustExec("INSERT INTO users(login, passwordhash, passwordsalt) VALUES ($1, $2, $3);", login, passwordHash, passwordSalt)
	tx.MustExec("INSERT INTO accrual(login, points, withdrawn) VALUES ($1, $2, $3);", login, decimal.Zero, decimal.Zero)

	err := tx.Commit()
	if err != nil {
		s.logger.Error("Error occurred while creating user", zap.Error(err))
		return err
	}

	return nil
}

// LoginUser gets user's data from the database to check for correct credentials
func (s storageLpq) LoginUser(login string) (string, string, error) {
	var (
		passwordHash string
		passwordSalt string
	)

	err := s.conn.QueryRow("SELECT passwordhash, passwordsalt FROM users WHERE login=$1", login).Scan(&passwordHash, &passwordSalt)
	if err != nil {
		s.logger.Error("Error occurred while authorizing user", zap.Error(err))
		return "", "", err
	}

	return passwordHash, passwordSalt, nil
}

// CheckDuplicateUser checks if user is already existing
func (s storageLpq) CheckDuplicateUser(login string) (bool, error) {
	var dbLogin string

	s.conn.QueryRow("SELECT login FROM users WHERE login=$1", login).Scan(&dbLogin)

	if dbLogin != "" {
		return true, nil
	}

	return false, nil
}
