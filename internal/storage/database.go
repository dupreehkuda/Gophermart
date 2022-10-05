package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type storage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(path string, logger *zap.Logger) *storage {
	conn, err := pgxpool.Connect(context.Background(), path)
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
    total      money,
    login     text not null
        unique
        references users,
    pointspaid integer,
    orderdate  date,
    accrual    integer,
    status     text
);`)
	batch.Queue(`
create table if not exists accrual
(
    login text    not null
        primary key
        references users,
    points integer not null
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
