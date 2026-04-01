package main

import (
	"context"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mrechkunov/gofermart/interal/compressmiddleware"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/handler"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/service"
	"github.com/mrechkunov/gofermart/server"
)

func main() {
	logger.Log.Infoln("reading config")
	config.Init()
	ctx := context.Background()
	r := chi.NewRouter()
	var Server = server.NewServer(&config.ConfigAddresses, r)
	go service.UpdateOrderListener(ctx, config.ChanToUpdate)
	r.Route("/api/user", func(r chi.Router) {
		r.Get("/orders", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersGet(ctx))))
		r.Get("/balance", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Balance(ctx))))
		r.Get("/withdrawals", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdrawals(ctx))))
		r.Post("/orders", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersPost(ctx, config.ChanToUpdate))))
		r.Post("/register", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Register(ctx))))
		r.Post("/login", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Login(ctx))))
		r.Post("/balance/withdraw", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdraw(ctx))))
	})

	logger.Log.Infoln("starting web server at:", config.ConfigAddresses.ServerBindAddress)
	Server.Run()
	close(config.ChanToUpdate)
	config.DBconn.Close()
}
