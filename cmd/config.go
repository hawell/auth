package main

import (
	"auth/database"
	"auth/logger"
	"auth/mailer"
	"auth/server"
)

type Config struct {
	Server   server.Config   `json:"server"`
	Mailer   mailer.Config   `json:"mailer"`
	Log      logger.Config   `json:"logger"`
	Database database.Config `json:"database"`
}
