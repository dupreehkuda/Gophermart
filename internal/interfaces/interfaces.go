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
	CheckUser(login string) (bool, error)
	CheckOrder(login string, order int) (bool, bool, error)
	NewOrder(login, status string, order int, accrual decimal.Decimal) error
	GetOrders(login string) ([]byte, error)
	GetBalance(login string) ([]byte, error)
	GetWithdrawals(login string) ([]byte, error)
	CheckPoints(order int, sum decimal.Decimal) (bool, error)
	WithdrawPoints(order int, sum decimal.Decimal) error
}

type Actions interface {
	Register(login, password string) (string, bool, error)
	Login(login, password string) (string, bool, error)
	NewOrder(login string, order int) error
	GetOrders(login string) ([]byte, error)
	GetBalance(login string) ([]byte, error)
	GetWithdrawals(login string) ([]byte, error)
	WithdrawPoints(order string, sum decimal.Decimal) error
}

type Middleware interface {
	CheckCompression(next http.Handler) http.Handler
	WriteCompressed(next http.Handler) http.Handler
	CheckToken(next http.Handler) http.Handler
}
