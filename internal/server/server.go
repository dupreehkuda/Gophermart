package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/dupreehkuda/Gophermart/internal/actions"
	"github.com/dupreehkuda/Gophermart/internal/configuration"
	"github.com/dupreehkuda/Gophermart/internal/handlers"
	intf "github.com/dupreehkuda/Gophermart/internal/interfaces"
	"github.com/dupreehkuda/Gophermart/internal/logger"
	"github.com/dupreehkuda/Gophermart/internal/middleware"
	"github.com/dupreehkuda/Gophermart/internal/storage"
)

type server struct {
	handlers intf.Handlers
	mw       intf.Middleware
	logger   *zap.Logger
	config   *configuration.Config
}

// NewByConfig creates new server instance with underlying layers
func NewByConfig() *server {
	log := logger.InitializeLogger()

	cfg := configuration.New(log)
	store := storage.New(cfg.DatabasePath, log)
	act := actions.New(store, log, cfg.AccrualAddress)
	handle := handlers.New(store, act, log)
	mware := middleware.New(act, log)

	return &server{
		handlers: handle,
		mw:       mware,
		logger:   log,
		config:   cfg,
	}
}

// Run runs the service
func (s server) Run() {
	serv := &http.Server{Addr: s.config.Address, Handler: s.service()}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				s.logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := serv.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Fatal("Error shutting down", zap.Error(err))
		}
		s.logger.Info("Server shut down", zap.String("port", s.config.Address))
		serverStopCtx()
	}()

	s.logger.Info("Server started", zap.String("port", s.config.Address))
	err := serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Cant start server", zap.Error(err))
	}

	<-serverCtx.Done()
}

func (s server) service() http.Handler {
	r := chi.NewRouter()

	r.Use(s.mw.CheckCompression)
	r.Use(s.mw.WriteCompressed)

	r.Route("/api/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", s.handlers.Register)
			r.Post("/login", s.handlers.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(s.mw.CheckToken)

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", s.handlers.AddOrder)
				r.Get("/", s.handlers.GetOrders)
			})

			r.Get("/withdrawals", s.handlers.GetWithdrawals)

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", s.handlers.GetBalance)
				r.Post("/withdraw", s.handlers.WithdrawPoints)
				r.Get("/withdrawals", s.handlers.GetWithdrawals)
			})
		})
	})

	return r
}
