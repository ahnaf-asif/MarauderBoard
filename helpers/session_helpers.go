package helpers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/google/uuid"
)

func SetSessionConfig() *session.Store {
	config := session.Config{
		Expiration: 24 * 30 * time.Hour,
		Storage: redis.New(redis.Config{
			Host:     "localhost",
			Port:     6379,
			Password: os.Getenv("REDIS_PASSWORD"),
			Username: "",
			Reset:    false,
		}),
		KeyLookup:      "cookie:_gothic_session",
		CookieDomain:   "",
		CookiePath:     "",
		CookieSecure:   os.Getenv("ENVIRONMENT") == "production",
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		KeyGenerator:   func() string { return uuid.New().String() },
	}
	sessionStore := session.New(config)
	return sessionStore
}
