package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading .env file")
		}
	}

	return nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
