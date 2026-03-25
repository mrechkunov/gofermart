package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/handler"
	"github.com/mrechkunov/gofermart/interal/logger"
)

func main() {
	config.Init()
	logger.Log.Infoln("reading config")
	r := chi.NewRouter()
	r.Post("/api/user/register", handler.Register)
	//	r.Post("/", handler.PostHandler)
	//	r.Get("/{id}", handler.GetHandler)
	logger.Log.Infoln("starting web server at:", config.ConfigAddresses.ServerBindAddress)
	err := http.ListenAndServe(config.ConfigAddresses.ServerBindAddress, r)
	if err != nil {
		logger.Log.Errorln(err)
	}

}
