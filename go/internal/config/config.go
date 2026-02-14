package config

import (
	"os"
)

// Config holds application configuration loaded from environment.
type Config struct {
	AppName   string
	Port      string
	DBURL     string
	Debug     bool
	Environment string
}

// Load reads configuration from environment variables.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/microservices"
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "go-microservice"
	}

	return &Config{
		AppName:     appName,
		Port:        port,
		DBURL:       dbURL,
		Debug:       os.Getenv("DEBUG") == "true",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
