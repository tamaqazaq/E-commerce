package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found. Using system environment variables")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
