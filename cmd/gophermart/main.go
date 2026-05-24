package main

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mrechkunov/gofermart/internal/config"
	"github.com/mrechkunov/gofermart/internal/logger"
	"github.com/mrechkunov/gofermart/internal/service"
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
	server.Run()
	close(config.ChanToUpdate)
	config.DBconn.Close()

	err := logger.Log.Sync()
	if err != nil && !errors.Is(err, syscall.EBADF) && !errors.Is(err, syscall.ENOTTY) {
		fmt.Println("error while zapLogger Sync in init function", err)
	}
}
