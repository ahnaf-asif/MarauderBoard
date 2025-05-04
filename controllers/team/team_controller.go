package team_controller

import "github.com/gofiber/fiber/v2"

func RegisterTeamController(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Team Home")
	})
}
