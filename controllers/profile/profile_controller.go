package profile_controller

import (
	"log"

	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RegisterProfileRoutes(app fiber.Router) {
	app.Get("/", middlewares.AuthMiddleware, func(ctx *fiber.Ctx) error {
		user, err := helpers.GetAuthUserSessionData(ctx)
		if err != nil {
			log.Println("User not found")
			return ctx.Redirect("/")
		}

		return ctx.Render("profile", fiber.Map{
			"User": user,
		}, "layouts/main")
	})
}
