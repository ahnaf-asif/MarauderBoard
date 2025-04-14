package controllers

import (
	"log"

	"github.com/ahnafasif/MarauderBoard/controllers/auth"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		user, err := helpers.GetAuthUserSessionData(ctx)
		if err != nil {
			log.Println("User not found")
		} else {
			db_user, _ := models.GetUserByEmail(database.DB, user.Email)
			return ctx.Render("index", fiber.Map{
				"User": db_user,
			}, "layouts/main")
		}
		return ctx.Render("index", fiber.Map{}, "layouts/main")
	})

	authGroup := app.Group("/auth")
	auth.RegisterAuthRoutes(authGroup)
}
