package main

import (
	"github.com/joho/godotenv"
	"github.com/meghashyamc/auth/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	godotenv.Load()
	cmd.Execute()

}
