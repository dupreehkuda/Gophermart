package service

import (
	"go.uber.org/zap"

	intf "github.com/dupreehkuda/Gophermart/internal/interfaces"
)

type Service struct {
	storage intf.Stored
	addr    string
	logger  *zap.Logger
	Channel chan int
	active  bool
}

func New(storage intf.Stored, logger *zap.Logger, addr string) *Service {
	ch := make(chan int, 10)
	serv := &Service{
		storage: storage,
		addr:    addr,
		logger:  logger,
		Channel: ch,
		active:  true,
	}

	if addr != "" {
		go serv.updateOrderData()
	}

	return serv
}