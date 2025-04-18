package auth

import (
	"encoding/json"
	"log"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func RegisterAuthRoutes(auth fiber.Router) {
	auth.Get("/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}
		return ctx.Redirect("/", fiber.StatusFound)
	})

	auth.Get(":provider", func(ctx *fiber.Ctx) error {
		_, err := goth_fiber.GetFromSession("user_data", ctx)
		if err != nil {
			return goth_fiber.BeginAuthHandler(ctx)
		}
		return ctx.Redirect("/", fiber.StatusFound)
	})

	auth.Get(":provider/callback", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			log.Println("Error completing user auth:", err)
		}

		found_user, err := models.GetUserByEmail(database.DB, user.Email)

		if err != nil {
			log.Println("User not found, creating new user:", err)
			new_user := models.User{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Provider:  &user.Provider,
				Email:     user.Email,
				Avatar:    &user.AvatarURL,
			}

			db_user, err := models.AddNewUser(database.DB, new_user)
			if err != nil {
				log.Fatal("AddNewUser Error happened here:", err)
			}
			found_user = db_user
		} else {
			found_user.FirstName = user.FirstName
			found_user.LastName = user.LastName
			found_user.Avatar = &user.AvatarURL
			if _, err := models.UpdateUser(database.DB, found_user); err != nil {
				log.Fatal(err)
			}
		}

		user_data := UserSessionData{
			ID:        found_user.ID,
			Email:     found_user.Email,
			FirstName: found_user.FirstName,
			LastName:  found_user.LastName,
			Provider:  *found_user.Provider,
			Avatar:    user.AvatarURL,
		}

		user_data_json, err := json.Marshal(user_data)
		if err != nil {
			log.Fatal(err)
		}

		if err := goth_fiber.StoreInSession("user_data", string(user_data_json), ctx); err != nil {
			log.Fatal(err)
		}

		redirect_uri := ctx.Cookies("redirect_uri", "/")
		return ctx.Redirect(redirect_uri, fiber.StatusFound)
	})
}
