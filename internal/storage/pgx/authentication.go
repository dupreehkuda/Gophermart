package pgx

import (
	"context"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CreateUser inserts new user's data in the database
func (s storage) CreateUser(login string, passwordHash string, passwordSalt string) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	pgxdecimal.Register(conn.Conn().TypeMap())

	batch := &pgx.Batch{}

	batch.Queue("INSERT INTO users(login, passwordhash, passwordsalt) VALUES ($1, $2, $3);", login, passwordHash, passwordSalt)
	batch.Queue("INSERT INTO accrual(login, points, withdrawn) VALUES ($1, $2, $3);", login, decimal.Zero, decimal.Zero)

	br := conn.SendBatch(context.Background(), batch)
	defer s.batchClosing(br)

	return nil
}

// LoginUser gets user's data from the database to check for correct credentials
func (s storage) LoginUser(login string) (string, string, error) {
	var (
		passwordHash string
		passwordSalt string
	)

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return "", "", err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(), "SELECT passwordhash, passwordsalt FROM users WHERE login=$1", login).Scan(&passwordHash, &passwordSalt)
	if err != nil {
		return "", "", err
	}

	return passwordHash, passwordSalt, nil
}

// CheckDuplicateUser checks if user is already existing
func (s storage) CheckDuplicateUser(login string) (bool, error) {
	var dbLogin string

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false, err
	}
	defer conn.Release()

	conn.QueryRow(context.Background(), "SELECT login FROM users WHERE login=$1", login).Scan(&dbLogin)

	if dbLogin != "" {
		return true, nil
	}

	return false, nil
}
