package dashboard_controller

import (
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterDashboardController(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		return ctx.Render("dashboard", data, "layouts/dashboard")
	})
}
