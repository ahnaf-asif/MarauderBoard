package project_controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	task_controller "github.com/ahnafasif/MarauderBoard/controllers/task"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterProjectControllers(app fiber.Router) {
	tasksGroup := app.Group("/:project_id/tasks")
	task_controller.RegisterTaskController(tasksGroup)

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
		todo_tasks, _ := models.GetTasksByProjectIdAndStatus(database.DB, uint(project_id), "Todo")
		in_progress_tasks, _ := models.GetTasksByProjectIdAndStatus(database.DB, uint(project_id), "In Progress")
		in_review_tasks, _ := models.GetTasksByProjectIdAndStatus(database.DB, uint(project_id), "In Review")
		done_tasks, _ := models.GetTasksByProjectIdAndStatus(database.DB, uint(project_id), "Done")
		cancelled_tasks, _ := models.GetTasksByProjectIdAndStatus(database.DB, uint(project_id), "Cancelled")

		tasks, _ := models.GetTasksByProjectId(database.DB, uint(project_id))
		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = project.Name
		data["Project"] = project
		data["Workspace"] = workspace
		data["TodoTasks"] = todo_tasks
		data["InProgressTasks"] = in_progress_tasks
		data["InReviewTasks"] = in_review_tasks
		data["DoneTasks"] = done_tasks
		data["CancelledTasks"] = cancelled_tasks
		data["Tasks"] = tasks

		return ctx.Render("projects/dashboard", data, "layouts/project")
	})

	app.Get("/:id/settings", func(ctx *fiber.Ctx) error {
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

		return ctx.Render("projects/settings", data, "layouts/project")
	})

	app.Get("/:id/backlog", func(ctx *fiber.Ctx) error {
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
		tasks, _ := models.GetTasksByProjectId(database.DB, uint(project_id))
		status_options := []string{"Todo", "In Progress", "In Review", "Done", "Cancelled"}
		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = project.Name
		data["Project"] = project
		data["Workspace"] = workspace
		data["Tasks"] = tasks
		data["StatusOptions"] = status_options

		// for all task, print task name, and then print names of task.Team.Users
		for _, task := range tasks {
			log.Println("Task Name:", task.Name)
			log.Println("Total Users: ", len(task.Team.Users))
		}
		return ctx.Render("projects/backlog", data, "layouts/project")
	})

	app.Get("/:id/kanban", func(ctx *fiber.Ctx) error {
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
		tasks, _ := models.GetTasksByProjectId(database.DB, uint(project_id))
		status_options := []string{"Todo", "In Progress", "In Review", "Done", "Cancelled"}
		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = project.Name
		data["Project"] = project
		data["Workspace"] = workspace
		data["Tasks"] = tasks
		data["StatusOptions"] = status_options

		return ctx.Render("projects/kanban", data, "layouts/project")
	})

	app.Get("/:id/gantt", func(ctx *fiber.Ctx) error {
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
		tasks, _ := models.GetTasksByProjectId(database.DB, uint(project_id))
		gantt_start_date := time.Now()
		gantt_end_date := time.Now().AddDate(0, 0, 14)

		for _, task := range tasks {
			if task.Status != "Cancelled" {
				if task.StartDate != nil && task.StartDate.Before(gantt_start_date) {
					gantt_start_date = *task.StartDate
				}
				if task.EndDate != nil && task.EndDate.After(gantt_end_date) {
					gantt_end_date = *task.EndDate
				}
			}

			switch task.Status {
			case "Todo":
				task.Progress = 0
			case "In Progress":
				task.Progress = 50
			case "In Review":
				task.Progress = 75
			case "Done":
				task.Progress = 100
			case "Cancelled":
				task.Progress = 0
			}
		}
		tasks_json, _ := json.Marshal(tasks)
		status_options := []string{"Todo", "In Progress", "In Review", "Done", "Cancelled"}
		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = project.Name
		data["Project"] = project
		data["Workspace"] = workspace
		data["Tasks"] = tasks
		data["StatusOptions"] = status_options
		data["TasksJson"] = string(tasks_json)
		data["StartDate"] = gantt_start_date
		data["EndDate"] = gantt_end_date

		return ctx.Render("projects/gantt", data, "layouts/project")
	})

	app.Post("/:id/update", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		project_id, _ := strconv.Atoi(id)

		project, err := models.GetProjectById(database.DB, uint(project_id))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}

		name := ctx.FormValue("name")
		description := ctx.FormValue("description")

		project.Name = name
		project.Description = description
		err = database.DB.Save(project).Error
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to update project",
			})
		}

		return ctx.Render("partials/success-message-with-disappear", fiber.Map{
			"Message": "Project updated successfully",
		})
	})

	app.Delete("/:id/delete", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		project_id, _ := strconv.Atoi(id)

		project, err := models.GetProjectById(database.DB, uint(project_id))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid project ID",
			})
		}

		err = models.DeleteProjectById(database.DB, uint(project_id))
		if err != nil {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"error": "Failed to delete project",
			})
		}
		return ctx.Render("partials/success-message-with-redirect", fiber.Map{
			"Message":  "Project deleted successfully",
			"Redirect": fmt.Sprintf("/workspaces/%d/projects", project.WorkspaceId),
		})
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
