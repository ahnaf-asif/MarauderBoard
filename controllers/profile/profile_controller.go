package profile_controller

import (
	"log"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/gofiber/fiber/v2"
)

func RegisterProfileRoutes(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		user, err := helpers.GetAuthUserSessionData(ctx)
		if err != nil {
			log.Println("User not found")
			return ctx.Redirect("/")
		}

		return ctx.Render("profile", fiber.Map{
			"User": user,
		}, "layouts/main")
	})

	app.Post("/update", func(ctx *fiber.Ctx) error {
		user, err := helpers.GetAuthUserSessionData(ctx)
		if err != nil {
			log.Println("User not found")
		}

		avatar, err := ctx.FormFile("avatar")
		var avatarURI string

		if err != nil {
			avatarURI = user.Avatar
		} else {
			avatarURI, err = helpers.UploadFile(avatar)
			if err != nil {
				log.Fatal("Error uploading file:", err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
			}
		}

		first_name := ctx.FormValue("first_name")
		last_name := ctx.FormValue("last_name")

		db_user, _ := models.GetUserByEmail(database.DB, user.Email)
		db_user.FirstName = first_name
		db_user.LastName = last_name
		if avatar != nil {
			db_user.Avatar = &avatarURI
		}

		db_user, err = models.UpdateUser(database.DB, db_user)
		if err != nil {
			log.Println("Error updating user:", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user",
			})
		}

		user.FirstName = db_user.FirstName
		user.LastName = db_user.LastName
		user.Avatar = *db_user.Avatar

		_, _ = helpers.SetAuthUserSessionData(ctx, user)

		return ctx.Render(
			"partials/success-message-with-disappear", fiber.Map{
				"Message":       "Profile Updated Successfully!",
				"RefreshAvatar": true,
				"Avatar":        user.Avatar,
			},
		)
	})
}
