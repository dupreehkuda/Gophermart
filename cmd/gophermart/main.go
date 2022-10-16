package main

import (
	"github.com/dupreehkuda/Gophermart/internal/configuration"
	"github.com/dupreehkuda/Gophermart/internal/handlers"
	"github.com/dupreehkuda/Gophermart/internal/logger"
	"github.com/dupreehkuda/Gophermart/internal/middleware"
	"github.com/dupreehkuda/Gophermart/internal/processors"
	"github.com/dupreehkuda/Gophermart/internal/server"
	"github.com/dupreehkuda/Gophermart/internal/storage"
)

func main() {
	log := logger.InitializeLogger()

	cfg := configuration.New(log)
	store := storage.New(cfg.DatabasePath, log)
	proc := processors.New(store, log, cfg.AccrualAddress)
	service := handlers.New(store, proc, log)
	mware := middleware.New(proc, log)

	api := server.New(service, mware, log)
	api.Launch(cfg.Address)
}
