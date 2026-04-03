package main

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/service"
	"github.com/mrechkunov/gofermart/server"
)

func main() {
	logger.Log.Infoln("reading config")
	config.Init()
	ctx := context.Background()
	var router = server.NewRouter(ctx)
	router.RoutesInit()
	server := server.NewServer(config.ConfigAddresses.ServerBindAddress, router.R)

	go service.UpdateOrderListener(ctx, config.ChanToUpdate)
	logger.Log.Infoln("starting web server at:", config.ConfigAddresses.ServerBindAddress)
	err := server.MyServer.ListenAndServe()
	if err != nil {
		logger.Log.Errorln("error while start http server")
	}
	close(config.ChanToUpdate)
	config.DBconn.Close()
}
