package main

import (
	"flag"
	"os"
)

type appConfig struct {
	listenAdd   string
	dbDsn       string
	accrualAddr string
	jwtKey      []byte
}

var _appConfig *appConfig

func init() {
	_appConfig = buildAppConfig()
}

func buildAppConfig() *appConfig {
	result := &appConfig{
		listenAdd:   "localhost:8090",
		accrualAddr: "localhost:8080",
		jwtKey:      []byte("super-secret-key"),
	}

	fillFromEnv(result)
	fillFromArgs(result)

	return result
}

func fillFromArgs(c *appConfig) {
	listenAdd := flag.String("a", "", "Net address host:port")
	dbDsn := flag.String("d", "", "Database DSN")
	accuralAddr := flag.String("r", "", "accural system address")
	flag.Parse()

	if listenAdd != nil {
		c.listenAdd = *listenAdd
	}
	if dbDsn != nil {
		c.dbDsn = *dbDsn
	}
	if accuralAddr != nil {
		c.listenAdd = *accuralAddr
	}
}

func fillFromEnv(c *appConfig) {
	if tmp, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		c.listenAdd = tmp
	}
	if tmp, exists := os.LookupEnv("DATABASE_URI"); exists {
		c.dbDsn = tmp
	}
	if tmp, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); exists {
		c.accrualAddr = tmp
	}
	if tmp, exists := os.LookupEnv("JWT_SECRET_KEY"); exists {
		c.jwtKey = []byte(tmp)
	}
}
