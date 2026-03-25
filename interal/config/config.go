package config

import (
	"flag"
	"os"
)

type Addresses struct {
	ServerBindAddress    string
	AccuralSystemAddress string
	MigrationsPath       string
	DBConnStr            string
}

var ConfigAddresses = Addresses{
	ServerBindAddress:    "localhost:8080",
	AccuralSystemAddress: "",
	DBConnStr:            "",
}

func Init() {
	ba := flag.String("a", "localhost:8080", "adress to server run")
	as := flag.String("r", "", "accural system address")
	mp := flag.String("m", "file://migrations", "default migration PATH")
	cs := flag.String("d", "postgres://yapra:yaprapass@10.254.40.123:5432/yandexpracticum?sslmode=disable", "default DBConnStr")
	flag.Parse()

	// если переиенные окружения установленны, берем их, иначе берем флаг
	if serverAddress, isEnvBindSrv := os.LookupEnv("RUN_ADDRESS"); isEnvBindSrv {
		ConfigAddresses.ServerBindAddress = serverAddress
	} else {
		ConfigAddresses.ServerBindAddress = *ba
	}

	if accuralAddress, isEnvAccural := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); isEnvAccural {
		ConfigAddresses.AccuralSystemAddress = accuralAddress
	} else {
		ConfigAddresses.AccuralSystemAddress = *as
	}

	if migratoinsPath, isEnvMigrationsPath := os.LookupEnv("MIGRATIONS_PATH"); isEnvMigrationsPath {
		ConfigAddresses.MigrationsPath = migratoinsPath
	} else {
		ConfigAddresses.MigrationsPath = *mp
	}

	if dbConnStr, isEnvDBConnStr := os.LookupEnv("DATABASE_URI"); isEnvDBConnStr {
		ConfigAddresses.DBConnStr = dbConnStr
	} else {
		ConfigAddresses.DBConnStr = *cs
	}
}
