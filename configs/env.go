package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port     string
	DbUrl    string
	RedisUrl string
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: ", err)
	}

	Port = getEnv("PORT", "42069")
	DbUrl = getEnv("DB_URL", "host=localhost user=ahnafasif password=postgres dbname=maruader_board port=5432")
	RedisUrl = getEnv("REDIS_URL", "localhost:6379")
}
