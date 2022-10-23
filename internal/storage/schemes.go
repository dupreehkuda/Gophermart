package storage

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type dbOrder struct {
	Number     string           `db:"orderID"`
	Status     string           `db:"status"`
	Accrual    decimal.Decimal  `db:"accrual"`
	UploadedAt pgtype.Timestamp `db:"orderdate"`
}

type order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type dbWithdrawal struct {
	Order       string
	Sum         decimal.Decimal
	ProcessedAt pgtype.Timestamp
}

type withdrawal struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type dbRespBalance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}

type respBalance struct {
	Current   float32
	Withdrawn float32
}
