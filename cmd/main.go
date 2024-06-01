package main

import (
	"auth/database"
	"auth/logger"
	"auth/mailer"
	"auth/server"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

func main() {
	configPtr := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	configFile := *configPtr
	var config Config
	err := cleanenv.ReadConfig(configFile, &config)
	if err != nil {
		panic(err)
	}

	eventLogger, err := logger.NewLogger(&config.Log.EventLog)
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(eventLogger)

	db, err := database.Connect(&config.Database)
	if err != nil {
		panic(err)
	}

	m, err := mailer.NewSMTP(&config.Mailer)
	if err != nil {
		panic(err)
	}

	accessLogger, err := logger.NewLogger(&config.Log.AccessLog)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	s := server.NewServer(&config.Server, db, m, accessLogger)
	err = s.ListenAndServer()
	fmt.Println(err)
}
