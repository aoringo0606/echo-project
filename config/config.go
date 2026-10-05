package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBPath    string
	AuthToken string
}

func Load() (Config, error) {
	authToken, ok := os.LookupEnv("AUTH_TOKEN")
	if !ok {
		return Config{}, fmt.Errorf("AUTH_TOKEN is required")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "users.db"
	}

	return Config{
		DBPath:    dbPath,
		AuthToken: authToken,
	}, nil
}