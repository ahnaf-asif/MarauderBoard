package workspace_controller

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/controllers/auth"
	"github.com/ahnafasif/MarauderBoard/controllers/project"
	team_controller "github.com/ahnafasif/MarauderBoard/controllers/team"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterWorkspaceControllers(app fiber.Router) {
	projectGroup := app.Group("/:workspace_id/projects")
	project_controller.RegisterProjectControllers(projectGroup)

	teamGroup := app.Group("/:workspace_id/teams")
	team_controller.RegisterTeamController(teamGroup)

	app.Get("/", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		user := data["User"].(auth.UserSessionData)

		workspaces, err := models.GetAllWorkspacesByUserId(database.DB, user.ID)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to fetch workspaces",
			})
		}

		data["PageTitle"] = "My Workspaces"
		data["Workspaces"] = workspaces
		return ctx.Render("workspaces/all-workspaces", data, "layouts/dashboard")
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

		ctx.Set("HX-Redirect", fmt.Sprintf("/workspaces/%d/dashboard", workspace.ID))
		return ctx.SendStatus(fiber.StatusNoContent)
	})

	app.Get("/:id/dashboard", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
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

		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = "Dashboard"
		data["Workspace"] = workspace

		return ctx.Render("workspaces/dashboard", data, "layouts/workspace")
	})

	app.Get("/:id/settings", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		workspace_id := ctx.Params("id")
		workspace_id_int, _ := strconv.Atoi(workspace_id)
		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id_int))
		data["Workspace"] = workspace
		all_users, _ := models.GetAllUsers(database.DB)
		available_users := []*models.User{}
		for _, user := range all_users {
			if user.ID != workspace.Administrator.ID {
				available_users = append(available_users, user)
			}
		}
		data["Available_Users"] = available_users
		return ctx.Render("workspaces/settings", data, "layouts/workspace")
	})

	app.Post("/:id/update", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		workspace_id, _ := strconv.Atoi(id)
		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id))
		name := ctx.FormValue("name")
		description := ctx.FormValue("description")

		workspace.Name = name
		workspace.Description = description
		err := database.DB.Save(workspace).Error
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to update workspace",
			})
		}
		return ctx.Render("partials/success-message-with-disappear", fiber.Map{
			"Message": "Workspace updated successfully",
		})
	})

	app.Post("/:id/change-admin", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		workspace_id, _ := strconv.Atoi(id)
		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id))
		new_admin_id := ctx.FormValue("new_admin_id")
		new_admin_id_int, _ := strconv.Atoi(new_admin_id)

		new_admin, _ := models.GetUserById(database.DB, uint(new_admin_id_int))
		workspace.AdministratorId = &new_admin.ID

		err := database.DB.Save(workspace).Error
		if err != nil {
			log.Println("Error updating workspace administrator:", err)
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to change workspace administrator",
			})
		}

		Notification := &models.Notification{
			Title:  "You are an administrator",
			UserId: new_admin.ID,
			Body:   "You have been assigned as the administrator of the workspace " + workspace.Name,
			Seen:   false,
		}
		_, _ = models.AddNotification(database.DB, Notification)

		return ctx.Render("partials/success-message-with-redirect", fiber.Map{
			"Message":  "Workspace updated successfully",
			"Redirect": "/workspaces",
		})
	})

	app.Delete("/:id/delete", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		workspace_id, _ := strconv.Atoi(id)
		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id))
		err := database.DB.Delete(workspace).Error
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to delete workspace",
			})
		}
		return ctx.Render("partials/success-message-with-redirect", fiber.Map{
			"Message":  "Workspace deleted successfully",
			"Redirect": "/workspaces",
		})
	})
}
