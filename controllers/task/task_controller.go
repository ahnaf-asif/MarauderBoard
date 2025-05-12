package task_controller

import (
	"html/template"
	"log"
	"strconv"

	comments_controller "github.com/ahnafasif/MarauderBoard/controllers/comment"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/ahnafasif/MarauderBoard/services"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterTaskController(app fiber.Router) {
	app.Get("/new", func(ctx *fiber.Ctx) error {
		project_id := ctx.Params("project_id")
		project_id_int, _ := strconv.Atoi(project_id)
		workspace_id := ctx.Params("workspace_id")
		workspace_id_int, _ := strconv.Atoi(workspace_id)
		project, _ := models.GetProjectById(database.DB, uint(project_id_int))
		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id_int))

		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = "Add Task"
		data["Project"] = project
		data["Workspace"] = workspace

		return ctx.Render("tasks/new", data, "layouts/task")
	})

	app.Post("/new", func(ctx *fiber.Ctx) error {
		project_id := ctx.Params("project_id")

		name := ctx.FormValue("name")
		description := ctx.FormValue("description")
		team_id := ctx.FormValue("team_id")

		project_id_int, _ := strconv.Atoi(project_id)
		team_id_int, _ := strconv.Atoi(team_id)
		project_id_uint := uint(project_id_int)
		team_id_uint := uint(team_id_int)

		task := &models.Task{
			Name:        name,
			Description: description,
			ProjectId:   &project_id_uint,
			TeamId:      &team_id_uint,
			Status:      "Todo",
		}
		redirect := "/workspaces/" + ctx.Params("workspace_id") + "/projects/" + project_id + "/backlog"

		db_task, _ := models.AddNewTask(database.DB, task)

		log.Println("Task created:", db_task.Name)

		return ctx.Render("partials/success-message-with-redirect", fiber.Map{
			"Message":  "Task created successfully",
			"Redirect": redirect,
		})
	})

	app.Post("/new/refine-description", func(ctx *fiber.Ctx) error {
		project_id := ctx.Params("project_id")
		workspace_id := ctx.Params("workspace_id")
		project_id_int, _ := strconv.Atoi(project_id)
		workspace_id_int, _ := strconv.Atoi(workspace_id)
		name := ctx.FormValue("name")
		description := ctx.FormValue("description")

		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id_int))
		project, _ := models.GetProjectById(database.DB, uint(project_id_int))

		context := "You are a project manager. You are given the name, descriptions of workspace and the project. You are also given a description. refine that description to make it more clear and conscise."
		context += "\n\nWorkspace Name: " + workspace.Name
		context += "\n\nWorkspace Description: " + workspace.Description
		context += "\n\nProject Name: " + project.Name
		context += "\n\nProject Description: " + project.Description
		context += "\n\nTask Name: " + name
		context += "\nStart with the title, and do not add anything extra in your message."
		context += "\nPut as much details and lists as possible"

		refined_description, _ := services.RefineText("gemma:2b", context, description)
		return ctx.Render("partials/tasks/description", fiber.Map{
			"Description": refined_description,
		})
	})

	app.Get("/:task_id/view", func(ctx *fiber.Ctx) error {
		project_id := ctx.Params("project_id")
		workspace_id := ctx.Params("workspace_id")
		task_id := ctx.Params("task_id")
		project_id_int, _ := strconv.Atoi(project_id)
		workspace_id_int, _ := strconv.Atoi(workspace_id)
		task_id_int, _ := strconv.Atoi(task_id)
		project, _ := models.GetProjectById(database.DB, uint(project_id_int))
		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id_int))
		task, _ := models.GetTaskById(database.DB, uint(task_id_int))

		rendered_description, _ := helpers.MarkdownWithMath(task.Description)

		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = "Task Details"
		data["Project"] = project
		data["Workspace"] = workspace
		data["Task"] = task
		data["Description"] = template.HTML(rendered_description)

		return ctx.Render("tasks/view", data, "layouts/task")
	})

	commentsController := app.Group("/:task_id/comments")
	comments_controller.RegisterCommentsController(commentsController)

	app.Post("/:task_id/status", func(ctx *fiber.Ctx) error {
		project_id := ctx.Params("project_id")
		workspace_id := ctx.Params("workspace_id")
		task_id := ctx.Params("task_id")
		task_id_int, _ := strconv.Atoi(task_id)
		status := ctx.FormValue("status")

		task, err := models.GetTaskById(database.DB, uint(task_id_int))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}

		prev_status := task.Status

		task.Status = status
		database.DB.Save(task)

		notification_link := "/workspaces/" + workspace_id + "/projects/" + project_id + "/backlog"

		if task.Assignee != nil {
			assignee_notification := &models.Notification{
				UserId: task.Assignee.ID,
				Title:  "Task status updated",
				Body:   "The status of task " + task.Name + " has been updated from " + prev_status + " to " + status,

				Seen: false,
				Link: &notification_link,
			}
			_, _ = models.AddNotification(database.DB, assignee_notification)

		}
		if task.Reporter != nil {
			reporter_notification := &models.Notification{
				UserId: task.Reporter.ID,
				Title:  "Task status updated",
				Body:   "The status of task " + task.Name + " has been updated from " + prev_status + " to " + status,
				Seen:   false,
				Link:   &notification_link,
			}
			_, _ = models.AddNotification(database.DB, reporter_notification)
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	})

	app.Post("/:task_id/assignee", func(ctx *fiber.Ctx) error {
		project_id := ctx.Params("project_id")
		workspace_id := ctx.Params("workspace_id")
		task_id := ctx.Params("task_id")
		task_id_int, _ := strconv.Atoi(task_id)
		assignee_id := ctx.FormValue("assignee_id")
		assignee_id_int, _ := strconv.Atoi(assignee_id)
		assignee_id_uint := uint(assignee_id_int)
		task, err := models.GetTaskById(database.DB, uint(task_id_int))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}
		// set task assignee to the new_assignee_uint
		assignee, _ := models.GetUserById(database.DB, assignee_id_uint)
		log.Println("New Assignee :", assignee.FirstName)

		previous_assignee := task.Assignee
		task.Assignee = assignee
		database.DB.Save(task)
		notification_link := "/workspaces/" + workspace_id + "/projects/" + project_id + "/backlog"
		// // send notification to the new assignee
		assignee_notification := &models.Notification{
			UserId: assignee.ID,
			Title:  "Task assigned to you",
			Body:   "You have been assigned to the task " + task.Name,
			Seen:   false,
			Link:   &notification_link,
		}
		_, _ = models.AddNotification(database.DB, assignee_notification)

		// send notification to Reporter
		if task.Reporter != nil {
			reporter_notification := &models.Notification{
				UserId: task.Reporter.ID,
				Title:  "Task assigned to " + task.Assignee.FirstName,
				Body:   "The task " + task.Name + " has been assigned to " + task.Assignee.FirstName,
				Seen:   false,
				Link:   &notification_link,
			}
			_, _ = models.AddNotification(database.DB, reporter_notification)

		}

		// send notification to previous Assignee
		if previous_assignee != nil {
			previous_assignee_notification := &models.Notification{
				UserId: previous_assignee.ID,
				Title:  "Task unassigned from you",
				Body:   "You have been unassigned from the task " + task.Name,
				Seen:   false,
				Link:   &notification_link,
			}
			_, _ = models.AddNotification(database.DB, previous_assignee_notification)
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	})

	app.Post("/:task_id/reporter", func(ctx *fiber.Ctx) error {
		project_id := ctx.Params("project_id")
		workspace_id := ctx.Params("workspace_id")
		task_id := ctx.Params("task_id")
		task_id_int, _ := strconv.Atoi(task_id)
		reporter_id := ctx.FormValue("reporter_id")
		reporter_id_int, _ := strconv.Atoi(reporter_id)
		reporter_id_uint := uint(reporter_id_int)
		task, err := models.GetTaskById(database.DB, uint(task_id_int))
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}
		reporter, _ := models.GetUserById(database.DB, reporter_id_uint)
		previous_reporter := task.Reporter

		task.Reporter = reporter
		database.DB.Save(task)

		notification_link := "/workspaces/" + workspace_id + "/projects/" + project_id + "/backlog"

		// send notification to the new assignee
		if task.Assignee != nil {
			assignee_notification := &models.Notification{
				UserId: task.Assignee.ID,
				Title:  "Task reporter updated",
				Body:   "The reporter of the task " + task.Name + " has been updated to " + reporter.FirstName,
				Seen:   false,
				Link:   &notification_link,
			}
			_, _ = models.AddNotification(database.DB, assignee_notification)
		}

		// send notification to reporter
		reporter_notification := &models.Notification{
			UserId: reporter.ID,
			Title:  "You are now the reporter of the task",
			Body:   "You are now the reporter of the task " + task.Name,
			Seen:   false,
			Link:   &notification_link,
		}
		_, _ = models.AddNotification(database.DB, reporter_notification)

		// send notification to previous Reporter
		if previous_reporter != nil {
			previous_reporter_notification := &models.Notification{
				UserId: previous_reporter.ID,
				Title:  "Task reporter updated",
				Body:   "The reporter of the task " + task.Name + " has been updated to " + reporter.FirstName,
				Seen:   false,
				Link:   &notification_link,
			}
			_, _ = models.AddNotification(database.DB, previous_reporter_notification)
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	})
}
