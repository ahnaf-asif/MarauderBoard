package project_controller

import (
	"fmt"
	"log"
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

	app.Get("/:id/dashboard", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		workspace_id := ctx.Params("workspace_id")
		workspace_id_int, _ := strconv.Atoi(workspace_id)
		project_id, _ := strconv.Atoi(id)

		project, err := models.GetProjectById(database.DB, uint(project_id))
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

	app.Get("/:id/teams", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		workspace_id := ctx.Params("workspace_id")
		workspace_id_int, _ := strconv.Atoi(workspace_id)
		project_id, _ := strconv.Atoi(id)

		project, err := models.GetProjectById(database.DB, uint(project_id))
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
		available_teams, err := models.GetAvailableTeamsForProject(database.DB, uint(workspace_id_int), uint(project_id))
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to fetch teams",
			})
		}

		log.Println("Available Teams(in project): ", available_teams)

		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = project.Name
		data["Project"] = project
		data["Workspace"] = workspace
		data["Available_Teams"] = available_teams

		return ctx.Render("projects/teams", data, "layouts/project")
	})

	app.Post("/:id/teams/add", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		team_id := ctx.FormValue("team_id")
		project_id, _ := strconv.Atoi(id)
		team_id_int, _ := strconv.Atoi(team_id)

		project, err := models.GetProjectById(database.DB, uint(project_id))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}

		err = models.AddTeamToProject(database.DB, uint(project_id), uint(team_id_int))
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to add team to project",
			})
		}

		ctx.Set("HX-Redirect", fmt.Sprintf("/workspaces/%d/projects/%d/teams", project.WorkspaceId, project.ID))
		return ctx.SendStatus(fiber.StatusNoContent)
	})

	app.Delete("/:id/teams/remove/:team_id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		team_id := ctx.Params("team_id")
		project_id, _ := strconv.Atoi(id)
		team_id_int, _ := strconv.Atoi(team_id)

		project, err := models.GetProjectById(database.DB, uint(project_id))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}

		err = models.RemoveTeamFromProject(database.DB, uint(project_id), uint(team_id_int))
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to remove team from project",
			})
		}

		ctx.Set("HX-Redirect", fmt.Sprintf("/workspaces/%d/projects/%d/teams", project.WorkspaceId, project.ID))
		return ctx.SendStatus(fiber.StatusNoContent)
	})
}
