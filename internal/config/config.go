package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		DatabaseURL string
		LogLevel    string
		AuddAPI
		LastFMAPI
	}

	AuddAPI struct {
		AuddAPIKey string
		AuddAPIURL string
	}

	LastFMAPI struct {
		LastFMAPIKey string
		LastFMAPIURL string
	}
)

func NewConfig() (*Config, error) {
	err := godotenv.Load("../../configs/.env")

	databaseURL := os.Getenv("DATABASE_URL")
	logLevel := os.Getenv("LOG_LEVEL")
	auddAPI := AuddAPI{os.Getenv("AUDD_API_KEY"), os.Getenv("AUDD_API_URL")}
	lastFMAPI := LastFMAPI{os.Getenv("LASTFM_API_KEY"), os.Getenv("LASTFM_API_URL")}

	if err != nil {
		return nil, err
	}

	conf := &Config{DatabaseURL: databaseURL, LogLevel: logLevel, AuddAPI: auddAPI, LastFMAPI: lastFMAPI}

	return conf, nil
}
