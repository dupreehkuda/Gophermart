package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	intf "github.com/dupreehkuda/Gophermart/internal/interfaces"
)

type server struct {
	handlers intf.Handlers
	mw       intf.Middleware
	logger   *zap.Logger
}

func New(handlers intf.Handlers, middleware intf.Middleware, logger *zap.Logger) *server {
	return &server{handlers: handlers, mw: middleware, logger: logger}
}

func (s server) Launch(address string) {
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

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", s.handlers.GetBalance)
				r.Post("/withdraw", s.handlers.WithdrawPoints)
				r.Get("/withdrawals", s.handlers.GetWithdrawals)
			})
		})
	})

	s.logger.Info("Server started", zap.String("port", address))
	err := http.ListenAndServe(address, r)
	if err != nil {
		s.logger.Error("Cant start server", zap.Error(err))
	}
}
