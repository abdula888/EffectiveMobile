package main

import (
	"EffectiveMobile/internal/app"
	"EffectiveMobile/internal/config"
	"log"
)

// @title           Music Library
// @version         1.0
// @description     API Server for Music Library.

// @host      localhost:8080
// @BasePath  /songs/
func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.Run(conf)
}
