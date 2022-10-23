package handlers

import (
	i "github.com/dupreehkuda/Gophermart/internal/interfaces"

	"go.uber.org/zap"
)

type handlers struct {
	storage i.Stored
	actions i.Actions
	logger  *zap.Logger
}

func New(storage i.Stored, processor i.Actions, logger *zap.Logger) *handlers {
	return &handlers{storage: storage, actions: processor, logger: logger}
}
