package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/gofermart/interal/compressmiddleware"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/handler"
	"github.com/mrechkunov/gofermart/interal/logger"
)

// Router - структура routera
type Router struct {
	ctx context.Context
	R   *chi.Mux
}

// NewServer Метод для инициализации структуры
func NewRouter(ctx context.Context) *Router {
	return &Router{
		ctx: ctx,
		R:   chi.NewRouter(),
	}
}

func (r *Router) RoutesInit() {
	r.R.Get("/api/user/orders", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersGet(r.ctx))))
	r.R.Get("/api/user/balance", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Balance(r.ctx))))
	r.R.Get("/api/user/withdrawals", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdrawals(r.ctx))))
	r.R.Post("/api/user/orders", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersPost(r.ctx, config.ChanToUpdate))))
	r.R.Post("/api/user/register", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Register(r.ctx))))
	r.R.Post("/api/user/login", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Login(r.ctx))))
	r.R.Post("/api/user/balance/withdraw", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdraw(r.ctx))))
}
