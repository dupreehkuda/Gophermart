package actions

import (
	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal/interfaces"
	s "github.com/dupreehkuda/Gophermart/internal/service"
)

type actions struct {
	service *s.Service
	storage i.Stored
	logger  *zap.Logger
	sysAddr string
}

func New(storage i.Stored, logger *zap.Logger, sysAddr string) *actions {
	serv := s.New(storage, logger, sysAddr)
	return &actions{service: serv,
		storage: storage,
		logger:  logger,
		sysAddr: sysAddr,
	}
}
