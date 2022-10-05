package handlers

import (
	i "github.com/dupreehkuda/Gophermart/internal/interfaces"
	"go.uber.org/zap"
)

type handlers struct {
	storage   i.Stored
	processor i.Processor
	logger    *zap.Logger
}

func New(storage i.Stored, processor i.Processor, logger *zap.Logger) *handlers {
	return &handlers{storage: storage, processor: processor, logger: logger}
}
