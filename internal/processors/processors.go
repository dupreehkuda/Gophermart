package processors

import (
	i "github.com/dupreehkuda/Gophermart/internal/interfaces"
	"go.uber.org/zap"
)

type processors struct {
	storage i.Stored
	logger  *zap.Logger
}

func New(storage i.Stored, logger *zap.Logger) *processors {
	return &processors{storage: storage, logger: logger}
}
