package actions

import (
	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal/interfaces"
)

type actions struct {
	storage i.Stored
	logger  *zap.Logger
	sysAddr string
}

func New(storage i.Stored, logger *zap.Logger, sysAddr string) *actions {
	return &actions{storage: storage, logger: logger, sysAddr: sysAddr}
}
