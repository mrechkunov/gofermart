package config

import (
	"flag"
)

type Adreses struct {
	ServerBindAdress   string
	ResultServerAdress string
}

var ConfigAdreses = Adreses{
	ServerBindAdress:   "localhost:8080",
	ResultServerAdress: "http://localhost:8080", // для работы unit теста
}

func Init() {
	ba := flag.String("a", "localhost:8080", "adress to server run")
	ra := flag.String("b", "http://localhost:8080", "default responce server adress")
	flag.Parse()
	ConfigAdreses.ServerBindAdress = *ba
	ConfigAdreses.ResultServerAdress = *ra
}
