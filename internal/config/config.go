package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file with err:", err.Error())
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
