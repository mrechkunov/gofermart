package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/logger"
)

// Server - структура сервера
type Server struct {
	AppConfig *config.Addresses
	Router    *chi.Mux
}

// NewServer Метод для инициализации структуры
func NewServer(appConfig *config.Addresses, router *chi.Mux) *Server {
	return &Server{
		AppConfig: appConfig,
		Router:    router,
	}
}

// Run - метод для запуска нашего http-сервера
func (server *Server) Run() {
	err := http.ListenAndServe(server.AppConfig.ServerBindAddress, server.Router)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
}
