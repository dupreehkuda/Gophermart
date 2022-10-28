package interfaces

import (
	"net/http"

	"github.com/shopspring/decimal"
)

type Handlers interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	AddOrder(w http.ResponseWriter, r *http.Request)
	GetOrders(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
	WithdrawPoints(w http.ResponseWriter, r *http.Request)
	GetWithdrawals(w http.ResponseWriter, r *http.Request)
}

type Stored interface {
	CreateUser(login, passwordHash, passwordSalt string) error
	LoginUser(login string) (string, string, error)
	CheckDuplicateUser(login string) (bool, error)
	CheckOrderExistence(login string, orderID int) (bool, bool)
	NewOrder(login string, orderID int) error
	GetOrders(login string) ([]byte, error)
	GetBalance(login string) ([]byte, error)
	GetWithdrawals(login string) ([]byte, error)
	CheckPoints(login string, sum decimal.Decimal) error
	WithdrawPoints(login string, orderID int, sum decimal.Decimal) error
	UpdateAccrual(orderID int, status string, accrual decimal.Decimal) error
}

type Actions interface {
	Register(login, password string) error
	Login(login, password string) error
	NewOrder(login string, orderID int) error
	GetOrders(login string) ([]byte, error)
	GetBalance(login string) ([]byte, error)
	GetWithdrawals(login string) ([]byte, error)
	WithdrawPoints(login, orderID string, sum decimal.Decimal) error
}

type Middleware interface {
	CheckCompression(next http.Handler) http.Handler
	WriteCompressed(next http.Handler) http.Handler
	CheckToken(next http.Handler) http.Handler
}
