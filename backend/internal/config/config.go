package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost 		string
	DBPort 		string
	DBUser 		string
	DBPassword 	string
	DBName 		string
	DBSSLMode 	string
}

func Load() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
        log.Println("No .env file found, using environment variables or defaults")
    }

	return &Config{
		DBHost:     getEnv("DB_HOST", ""),
        DBPort:     getEnv("DB_PORT", ""),
        DBUser:     getEnv("DB_USER", ""),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", ""),
        DBSSLMode:  getEnv("DB_SSL_MODE", ""),
	}
}

func (c *Config) GetDBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}