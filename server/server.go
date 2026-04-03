package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/logger"
)

// Server - структура сервера
type ServerHTTP struct {
	MyServer http.Server
	Router   *chi.Mux
}

// NewServer Метод для инициализации структуры
func NewServer(address string, router *chi.Mux) *ServerHTTP {
	var result ServerHTTP
	result.MyServer.Addr = address
	result.Router = router
	return &result
}

// Run - метод для запуска нашего http-сервера
func (s *ServerHTTP) Run() {

	err := http.ListenAndServe(config.ConfigAddresses.ServerBindAddress, s.Router)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
}
