package controllers

import (
	"log"

	ai_controller "github.com/ahnafasif/MarauderBoard/controllers/ai"
	"github.com/ahnafasif/MarauderBoard/controllers/auth"
	"github.com/ahnafasif/MarauderBoard/controllers/workspace"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/middlewares"
	"github.com/gofiber/fiber/v2"
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

	app.Get("/profile", middlewares.AuthMiddleware, func(ctx *fiber.Ctx) error {
		user, err := helpers.GetAuthUserSessionData(ctx)
		if err != nil {
			log.Println("User not found")
			return ctx.Redirect("/")
		}
		return ctx.Render("profile", fiber.Map{
			"User": user,
		}, "layouts/main")
	})

	app.Get("/dashboard", middlewares.AuthMiddleware, func(ctx *fiber.Ctx) error {
		user := ctx.Locals("User")
		return ctx.Render("dashboard", fiber.Map{
			"User": user,
		}, "layouts/dashboard")
	})

	authGroup := app.Group("/auth")
	auth.RegisterAuthRoutes(authGroup)

	workspaceGroup := app.Group("/workspaces", middlewares.AuthMiddleware)
	workspace_controller.RegisterWorkspaceControllers(workspaceGroup)

	aiGroup := app.Group("/ai", middlewares.AuthMiddleware)
	ai_controller.RegisterAiControllers(aiGroup)
}
