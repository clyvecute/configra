package config

import (
	"os"

	"github.com/clyvecute/configra/internal/db"
)

type AppConfig struct {
	DB   db.Config
	Port string
}

func Load() AppConfig {
	return AppConfig{
		Port: getEnv("PORT", "8080"),
		DB: db.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "configra"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
