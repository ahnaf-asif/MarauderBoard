package team_controller

import (
	"log"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterTeamController(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		workspace_id, _ := ctx.ParamsInt("workspace_id")

		workspace, err := models.GetWorkspaceById(database.DB, uint(workspace_id))
		if err != nil {
			return ctx.Status(404).SendString("Workspace not found")
		}

		data["PageName"] = "Teams"
		data["Workspace"] = workspace

		teams, err := models.GetTeamsByWorkspaceId(database.DB, workspace.ID)
		if err != nil {
			return ctx.Status(500).SendString("Error getting teams")
		}

		data["Teams"] = teams

		return ctx.Render("workspace/teams", data, "layouts/workspace")
	})

	app.Get("/new", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		workspace_id, _ := ctx.ParamsInt("workspace_id")

		workspace, err := models.GetWorkspaceById(database.DB, uint(workspace_id))
		if err != nil {
			return ctx.Status(404).SendString("Workspace not found")
		}

		data["PageName"] = "New Team"
		data["Workspace"] = workspace

		all_users, err := models.GetAllUsers(database.DB)
		if err != nil {
			return ctx.Status(500).SendString("Error getting users")
		}

		data["Users"] = &all_users

		return ctx.Render("teams/create-team", data, "layouts/workspace")
	})

	app.Post("/create", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
		leader_id := ctx.FormValue("leader_id")
		leader_id_int, _ := strconv.Atoi(leader_id)
		workspace_id, _ := ctx.ParamsInt("workspace_id")

		log.Println(name, leader_id, workspace_id)

		chat_group := models.ChatGroup{
			Name: name,
		}
		_, _ = models.AddChatGroup(database.DB, &chat_group)

		new_team := models.Team{
			Name:        name,
			LeaderId:    uint(leader_id_int),
			WorkspaceId: uint(workspace_id),
			ChatGroupId: chat_group.ID,
		}
		_, err := models.AddNewTeam(database.DB, &new_team)
		if err != nil {
			return ctx.Status(500).SendString("Error creating team")
		}

		leader, err := models.GetUserById(database.DB, uint(leader_id_int))
		if err != nil {
			return ctx.Status(500).SendString("Error getting leader")
		}

		notification_link := "/workspaces/" + strconv.Itoa(workspace_id) + "/teams/" + strconv.Itoa(int(new_team.ID))
		leader_notification := models.Notification{
			Title:  "New Team Created",
			Body:   "You have been assigned as the leader of the team " + new_team.Name,
			UserId: leader.ID,
			Seen:   false,
			Link:   &notification_link,
		}

		_, err = models.AddNotification(database.DB, &leader_notification)
		if err != nil {
			return ctx.Status(500).SendString("Error creating notification")
		}

		new_team.Users = append(new_team.Users, leader)
		_, err = models.UpdateTeam(database.DB, &new_team)
		if err != nil {
			return ctx.Status(500).SendString("Error updating team")
		}

		return ctx.Render("partials/success-message-with-redirect", fiber.Map{
			"Message":  "Team created successfully.",
			"Redirect": "/workspaces/" + strconv.Itoa(workspace_id) + "/teams",
		})
	})
}
