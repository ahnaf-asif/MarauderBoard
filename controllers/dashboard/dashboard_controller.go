package dashboard_controller

import (
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/gofiber/fiber/v2"
)

func RegisterDashboardController(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		user, _ := helpers.GetAuthUserSessionData(ctx)
		return ctx.Render("dashboard", fiber.Map{
			"User": user,
		}, "layouts/dashboard")
	})
}
