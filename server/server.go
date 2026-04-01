package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/logger"
)

// Server - структура сервера
type ServerHTTP struct {
	Addr   string
	Router *chi.Mux
}

// NewServer Метод для инициализации структуры
func NewServer(appConfig *config.Addresses, router *chi.Mux) *ServerHTTP {
	return &ServerHTTP{
		Addr:   config.ConfigAddresses.ServerBindAddress,
		Router: router,
	}
}

// Run - метод для запуска нашего http-сервера
func (s *ServerHTTP) Run() {
	err := http.ListenAndServe(config.ConfigAddresses.ServerBindAddress, s.Router)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
}
