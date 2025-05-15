package dashboard_controller

import (
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterDashboardController(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		user, _ := helpers.GetAuthUserSessionData(ctx)

		workspaces, _ := models.GetAllWorkspacesByUserId(database.DB, user.ID)
		data["Workspaces"] = workspaces

		return ctx.Render("dashboard", data, "layouts/dashboard")
	})
}
