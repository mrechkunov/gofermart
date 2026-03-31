package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mrechkunov/gofermart/interal/compressmiddleware"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/handler"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/service"
)

func main() {
	config.Init()
	ctx := context.Background()
	logger.Log.Infoln("reading config")
	r := chi.NewRouter()
	chanToUpdate := make(chan int64, 10)
	go service.UpdateOrderListener(ctx, chanToUpdate)
	r.Route("/api/user", func(r chi.Router) {
		r.Get("/orders/", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersGet)))
		r.Get("/balance", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Balance)))
		r.Get("/withdrawals", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdrawals)))
		r.Post("/orders/", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersPost(chanToUpdate))))
		r.Post("/register", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Register)))
		r.Post("/login", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Login)))
		r.Post("/balance/withdraw", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdraw(ctx))))
	})
	logger.Log.Infoln("starting web server at:", config.ConfigAddresses.ServerBindAddress)
	err := http.ListenAndServe(config.ConfigAddresses.ServerBindAddress, r)
	if err != nil {
		logger.Log.Errorln(err)
	}
	close(chanToUpdate)
	config.DBconn.Close()
}
