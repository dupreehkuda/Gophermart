package interfaces

import "net/http"

type Handlers interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	AddOrder(w http.ResponseWriter, r *http.Request)
	GetOrders(w http.ResponseWriter, r *http.Request)
}

type Stored interface {
	CreateUser(login, passwordHash, passwordSalt string) error
	LoginUser(login string) (string, string, error)
	CheckUser(login string) (bool, error)
	CheckOrder(login string, order int) (bool, bool, error)
	NewOrder(login, status string, order, accrual int) error
	GetOrders(login string) ([]byte, error)
}

type Processor interface {
	Register(login, password string) (string, bool, error)
	Login(login, password string) (string, bool, error)
	NewOrder(login string, order int) (bool, bool, bool, error)
	GetOrders(login string) ([]byte, error)
}

type Middleware interface {
	CheckCompression(next http.Handler) http.Handler
	WriteCompressed(next http.Handler) http.Handler
	CheckToken(next http.Handler) http.Handler
}
