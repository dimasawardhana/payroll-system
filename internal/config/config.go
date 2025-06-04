package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl      string
	JWTSecret  string
	ServerPort string
	Env        string // add Env for environment
}

func Load() *Config {
	_ = godotenv.Load() // Load .env file if present
	return &Config{
		DBUrl:      os.Getenv("DATABASE_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("SERVER_PORT"),
		Env:        os.Getenv("ENV"), // load ENV from environment
	}
}
