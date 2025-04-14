package helpers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory" // In-memory storage
	"github.com/google/uuid"
)

func SetSessionConfig() *session.Store {
	config := session.Config{
		Expiration:     24 * 30 * time.Hour,
		Storage:        memory.New(),
		KeyLookup:      "cookie:_gothic_session",
		CookieDomain:   "",
		CookiePath:     "",
		CookieSecure:   os.Getenv("ENVIRONMENT") == "production",
		CookieHTTPOnly: true, // Should always be enabled
		CookieSameSite: "Lax",
		KeyGenerator:   func() string { return uuid.New().String() },
	}
	sessionStore := session.New(config)
	return sessionStore
}
