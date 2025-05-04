package controllers

import (
	"log"

	ai_controller "github.com/ahnafasif/MarauderBoard/controllers/ai"
	"github.com/ahnafasif/MarauderBoard/controllers/auth"
	dashboard_controller "github.com/ahnafasif/MarauderBoard/controllers/dashboard"
	notifications_controller "github.com/ahnafasif/MarauderBoard/controllers/notifications"
	profile_controller "github.com/ahnafasif/MarauderBoard/controllers/profile"
	"github.com/ahnafasif/MarauderBoard/controllers/workspace"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		user, err := helpers.GetAuthUserSessionData(ctx)
		if err != nil {
			log.Println("User not found")
		} else {
			return ctx.Render("index", fiber.Map{
				"User": user,
			}, "layouts/main")
		}
		return ctx.Render("index", fiber.Map{}, "layouts/main")
	})

	app.Get("/auth/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}
		return ctx.Redirect("/", fiber.StatusFound)
	})

	authGroup := app.Group("/auth", middlewares.GuestMiddleware)
	auth.RegisterAuthRoutes(authGroup)

	profileGroup := app.Group("/profile", middlewares.AuthMiddleware)
	profile_controller.RegisterProfileRoutes(profileGroup)

	dashboardGroup := app.Group("/dashboard", middlewares.AuthMiddleware)
	dashboard_controller.RegisterDashboardController(dashboardGroup)

	workspaceGroup := app.Group("/workspaces", middlewares.AuthMiddleware)
	workspace_controller.RegisterWorkspaceControllers(workspaceGroup)

	aiGroup := app.Group("/ai", middlewares.AuthMiddleware)
	ai_controller.RegisterAiControllers(aiGroup)

	notificationsGroup := app.Group("/notifications", middlewares.AuthMiddleware)
	notifications_controller.RegisterNotificationController(notificationsGroup)
}
