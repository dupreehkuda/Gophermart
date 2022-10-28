package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type storage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(path string, logger *zap.Logger) *storage {
	config, err := pgxpool.ParseConfig(path)
	if err != nil {
		logger.Error("Unable to parse config", zap.Error(err))
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("Unable to connect to database", zap.Error(err))
	}

	batch := &pgx.Batch{}
	batch.Queue(`
create table if not exists users
(
    login        text not null
        primary key unique,
    passwordhash text not null,
    passwordsalt text not null
        unique
);`)
	batch.Queue(`
create table if not exists orders
(
    orderid    bigint not null
        unique,
    login     text not null
        references users,
    pointsspent bool,
    orderdate  timestamp,
    accrual    numeric,
    status     text
);`)
	batch.Queue(`
create table if not exists accrual
(
    login text    not null
        primary key
        references users,
    points numeric not null,
    withdrawn numeric
);
`)

	br := conn.SendBatch(context.Background(), batch)
	err = br.Close()
	if err != nil {
		logger.Error("Error occurred while creating table", zap.Error(err))
	}

	return &storage{
		pool:   conn,
		logger: logger,
	}
}

func (s storage) batchClosing(br pgx.BatchResults) {
	err := br.Close()
	if err != nil {
		s.logger.Error("Error while closing batch", zap.Error(err))
		return
	}
}
