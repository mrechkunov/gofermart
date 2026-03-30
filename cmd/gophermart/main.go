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
	r.Route("/api/user/orders", func(r chi.Router) {
		r.Get("/", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersGet)))
		r.Post("/", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.OrdersPost(chanToUpdate))))
	})
	r.Post("/api/user/register", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Register)))
	r.Post("/api/user/login", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Login)))
	r.Get("/api/user/balance", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Balance)))
	r.Post("/api/user/balance/withdraw", logger.WithLogging(compressmiddleware.GzipMiddleware(handler.Withdraw)))

	logger.Log.Infoln("starting web server at:", config.ConfigAddresses.ServerBindAddress)
	err := http.ListenAndServe(config.ConfigAddresses.ServerBindAddress, r)
	if err != nil {
		logger.Log.Errorln(err)
	}
	close(chanToUpdate)
	config.DBconn.Close()
}
