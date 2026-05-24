package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/gofermart/internal/compressmiddleware"
	"github.com/mrechkunov/gofermart/internal/config"
	"github.com/mrechkunov/gofermart/internal/cryptoauth"
	"github.com/mrechkunov/gofermart/internal/handler"
	"github.com/mrechkunov/gofermart/internal/logger"
)

// Router - структура routera
type Router struct {
	R *chi.Mux
}

// NewServer Метод для инициализации структуры
func NewRouter(ctx context.Context) *Router {
	return &Router{R: chi.NewRouter()}
}

func (r *Router) RoutesInit() {
	r.R.Get("/api/user/orders", cryptoauth.WithAuth(logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersGet))))
	r.R.Get("/api/user/balance", cryptoauth.WithAuth(logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Balance))))
	r.R.Get("/api/user/withdrawals", cryptoauth.WithAuth(logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdrawals))))
	r.R.Post("/api/user/orders", cryptoauth.WithAuth(logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersPost(config.ChanToUpdate)))))
	r.R.Post("/api/user/register", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Register)))
	r.R.Post("/api/user/login", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Login)))
	r.R.Post("/api/user/balance/withdraw", cryptoauth.WithAuth(logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdraw))))
}
