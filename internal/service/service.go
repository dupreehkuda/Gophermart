package service

import (
	"go.uber.org/zap"

	intf "github.com/dupreehkuda/Gophermart/internal/interfaces"
)

type Service struct {
	storage    intf.Stored
	addr       string
	logger     *zap.Logger
	OrderQueue chan int
	active     bool
}

// New creates new instance of accrual call system
func New(storage intf.Stored, logger *zap.Logger, addr string) *Service {
	ch := make(chan int, 10)
	serv := &Service{
		storage:    storage,
		addr:       addr,
		logger:     logger,
		OrderQueue: ch,
		active:     true,
	}

	if addr != "" {
		go serv.updateOrderData()
	}

	return serv
}
