package workspace_controller

import (
	"fmt"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/gofiber/fiber/v2"
)

func RegisterWorkspaceControllers(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		user, _ := helpers.GetAuthUserSessionData(ctx)
		workspaces, err := models.GetAllWorkspacesByUserId(database.DB, user.ID)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to fetch workspaces",
			})
		}
		return ctx.Render("workspaces/all-workspaces", fiber.Map{
			"PageTitle":  "My Workspaces",
			"Workspaces": workspaces,
			"User":       user,
		}, "layouts/dashboard")
	})

	app.Get("/get-started", func(ctx *fiber.Ctx) error {
		return ctx.Render("workspaces/get-started", fiber.Map{})
	})

	app.Post("/create", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
		description := ctx.FormValue("description")
		user, _ := helpers.GetAuthUserSessionData(ctx)

		chatGroup := &models.ChatGroup{
			Name: name,
		}

		chatGroup, err := models.AddChatGroup(database.DB, chatGroup)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to create chat group",
			})
		}

		workspace := &models.Workspace{
			Name:            name,
			AdministratorId: &user.ID,
			Description:     description,
			ChatGroupId:     &chatGroup.ID,
		}
		workspace, err = models.AddWorkspace(database.DB, workspace)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to create workspace",
			})
		}

		ctx.Set("HX-Redirect", fmt.Sprintf("/workspaces/%d", workspace.ID))
		return ctx.SendStatus(fiber.StatusNoContent)
	})

	app.Get("/:id/dashboard", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		user, _ := helpers.GetAuthUserSessionData(ctx)
		if id == "" {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Workspace ID is required",
			})
		}

		workspaceId, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid workspace ID",
			})
		}

		workspace, err := models.GetWorkspaceById(database.DB, uint(workspaceId))
		if err != nil {
			return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
				"error": "Workspace not found",
			})
		}

		return ctx.Render("workspaces/dashboard", fiber.Map{
			"PageTitle": "Dashboard",
			"Workspace": workspace,
			"User":      user,
		}, "layouts/workspace")
	})
}
