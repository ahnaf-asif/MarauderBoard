package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Host               string
	Port               string
	DbUrl              string
	RedisUrl           string
	AuthKey            string
	AuthMaxAge         int
	Environment        string
	GoogleClientId     string
	GoogleClientSecret string
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

	Host = getEnv("HOST", "127.0.0.1")
	Port = getEnv("PORT", "42069")
	Environment = getEnv("ENVIRONMENT", "development")
	AuthMaxAge, _ = strconv.Atoi(getEnv("AUTH_MAX_AGE", "3600"))
	AuthKey = getEnv("AUTH_KEY", "secret")
	DbUrl = getEnv("DB_URL", "host=localhost user=ahnafasif password=postgres dbname=maruader_board port=5432")
	RedisUrl = getEnv("REDIS_URL", "localhost:6379")
	GoogleClientId = getEnv("GOOGLE_CLIENT_ID", "")
	GoogleClientSecret = getEnv("GOOGLE_CLIENT_SECRET", "")
}
