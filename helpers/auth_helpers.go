package helpers

import (
	"encoding/json"
	"log"

	"github.com/ahnafasif/MarauderBoard/configs"
	"github.com/ahnafasif/MarauderBoard/controllers/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

func SetAuthProviders() {
	goth.UseProviders(
		google.New(
			configs.GoogleClientId,
			configs.GoogleClientSecret,
			"http://127.0.0.1:42069/auth/google/callback",
			"email",
			"profile",
		),
	)
}

func GetAuthUserSessionData(ctx *fiber.Ctx) (auth.UserSessionData, error) {
	var user_data auth.UserSessionData

	// Get JSON string from session
	user_data_json, err := goth_fiber.GetFromSession("user_data", ctx)
	if err != nil {
		return user_data, err
	}

	// Unmarshal JSON
	if err := json.Unmarshal([]byte(user_data_json), &user_data); err != nil {
		return user_data, err
	}

	log.Println("User data in goth session: ", user_data)

	return user_data, nil
}
