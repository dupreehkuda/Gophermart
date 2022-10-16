package processors

import (
	i "github.com/dupreehkuda/Gophermart/internal/interfaces"
	"go.uber.org/zap"
)

type processors struct {
	storage i.Stored
	logger  *zap.Logger
	sysAddr string
}

func New(storage i.Stored, logger *zap.Logger, sysAddr string) *processors {
	return &processors{storage: storage, logger: logger, sysAddr: sysAddr}
}
