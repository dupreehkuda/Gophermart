package middleware

import (
	"go.uber.org/zap"

	i "github.com/dupreehkuda/Gophermart/internal/interfaces"
)

type middleware struct {
	processor i.Processor
	logger    *zap.Logger
}

func New(processor i.Processor, logger *zap.Logger) *middleware {
	return &middleware{processor: processor, logger: logger}
}
