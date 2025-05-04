package project_controller

import (
	"fmt"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
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

		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = "All Projects"
		data["Workspace"] = workspace
		data["Projects"] = projects

		return ctx.Render("projects/all-projects", data, "layouts/workspace")
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

		data := load_locals.LoadLocals(ctx)
		data["Workspace"] = workspace
		return ctx.Render("projects/create", data)
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

		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = project.Name
		data["Project"] = project
		data["Workspace"] = workspace

		return ctx.Render("projects/dashboard", data, "layouts/project")
	})
}
