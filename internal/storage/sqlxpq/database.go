package sqlxpq

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type storageLpq struct {
	conn   *sqlx.DB
	logger *zap.Logger
}

var schema = `
create table if not exists users (
    login        text not null primary key unique,
    passwordhash text not null,
    passwordsalt text not null unique
);

create table if not exists orders (
    orderid    bigint not null primary key unique,
    login     text not null references users,
    pointsspent bool default FALSE,
    orderdate  timestamp,
    accrual    numeric,
    status     text
);

create table if not exists accrual (
    login text not null primary key references users,
    points numeric not null,
    withdrawn numeric not null
);

create table if not exists withdrawals (
    orderid bigint not null unique,
    login text not null primary key references users, 
    withdrawn numeric not null,
    processed_at timestamp
);
`

// New creates a new instance of database layer and migrates it
func New(path string, logger *zap.Logger) *storageLpq {
	db, err := sqlx.Connect("postgres", path)
	if err != nil {
		logger.Error("Unable to connect db", zap.Error(err))
	}

	db.MustExec(schema)

	logger.Info("Launched with sqlx")

	return &storageLpq{
		conn:   db,
		logger: logger,
	}
}
