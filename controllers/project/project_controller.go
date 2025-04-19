package project_controller

import (
	"fmt"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/gofiber/fiber/v2"
)

func RegisterProjectControllers(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		workspaceId := ctx.Params("workspace_id")
		workspaceIdInt, _ := strconv.Atoi(workspaceId)

		workspace, err := models.GetWorkspaceById(database.DB, uint(workspaceIdInt))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid workspace ID",
			})
		}

		projects, err := models.GetProjectsByWorkspaceId(database.DB, workspaceIdInt)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to fetch projects",
			})
		}

		user, _ := helpers.GetAuthUserSessionData(ctx)
		return ctx.Render("projects/all-projects", fiber.Map{
			"PageTitle": "All Projects",
			"User":      user,
			"Workspace": workspace,
			"Projects":  projects,
		}, "layouts/workspace")
	})

	app.Get("/create", func(ctx *fiber.Ctx) error {
		workspaceId := ctx.Params("workspace_id")
		workspaceIdInt, _ := strconv.Atoi(workspaceId)

		workspace, err := models.GetWorkspaceById(database.DB, uint(workspaceIdInt))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid workspace ID",
			})
		}

		user, _ := helpers.GetAuthUserSessionData(ctx)
		return ctx.Render("projects/create", fiber.Map{
			"User":      user,
			"Workspace": workspace,
		})
	})

	app.Post("/create", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
		description := ctx.FormValue("description")
		workspaceId := ctx.Params("workspace_id")
		workspaceIdInt, _ := strconv.Atoi(workspaceId)

		project := &models.Project{
			Name:        name,
			Description: description,
			Status:      "Pending",
			WorkspaceId: workspaceIdInt,
		}

		_, err := models.AddNewProject(database.DB, project)
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to create project",
			})
		}
		ctx.Set("HX-Redirect", fmt.Sprintf("/workspaces/%d/projects", workspaceIdInt))
		return ctx.SendStatus(fiber.StatusNoContent)
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		workspace_id := ctx.Params("workspace_id")
		workspace_id_int, _ := strconv.Atoi(workspace_id)
		project_id, _ := strconv.Atoi(id)

		project, err := models.GetProjectById(database.DB, project_id)
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}

		workspace, err := models.GetWorkspaceById(database.DB, uint(workspace_id_int))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid workspace ID",
			})
		}

		user, _ := helpers.GetAuthUserSessionData(ctx)
		return ctx.Render("projects/dashboard", fiber.Map{
			"User":      user,
			"Project":   project,
			"Workspace": workspace,
		}, "layouts/project")
	})
}
