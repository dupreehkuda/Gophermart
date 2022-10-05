package storage

import (
	"context"
	"go.uber.org/zap"
)

func (s storage) CheckOrder(login string, order int) (bool, bool, error) {
	var dbUserLogin string

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false, false, err
	}
	defer conn.Release()

	conn.QueryRow(context.Background(), "select login from orders where orderid = $1;", order).Scan(&dbUserLogin)

	if dbUserLogin == "" {
		return false, false, nil
	}

	if dbUserLogin != login {
		return true, false, nil
	}

	return true, true, nil
}

func (s storage) NewOrder(login string, order int) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	_, err = conn.Query(context.Background(), "insert into orders(login, orderid) values ($1, $2);", login, order)
	if err != nil {
		s.logger.Error("Error while inserting order", zap.Error(err))
		return err
	}

	return nil
}
