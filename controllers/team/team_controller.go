package team_controller

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
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
		for _, team := range teams {
			log.Println("Team: ", team.Leader.Email, team.Leader.ID)
		}
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
		user, _ := helpers.GetAuthUserSessionData(ctx)
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

		leader, err := models.GetUserById(database.DB, uint(leader_id_int))
		if err != nil {
			return ctx.Status(500).SendString("Error getting leader")
		}
		if leader.ID != user.ID {
			notification_link := "/workspaces/" +
				strconv.Itoa(workspace_id) + "/teams/" +
				strconv.Itoa(int(new_team.ID))
			leader_notification := models.Notification{
				Title: "New Team Created",
				Body: "You have been assigned as the leader of the team " +
					new_team.Name,
				UserId: leader.ID,
				Seen:   false,
				Link:   &notification_link,
			}

			_, err = models.AddNotification(database.DB, &leader_notification)
			if err != nil {
				return ctx.Status(500).SendString("Error creating notification")
			}
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

	app.Get("/:team_id", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		workspace_id, _ := ctx.ParamsInt("workspace_id")
		team_id, _ := ctx.ParamsInt("team_id")
		workspace, err := models.GetWorkspaceById(database.DB, uint(workspace_id))
		if err != nil {
			return ctx.Status(404).SendString("Workspace not found")
		}
		data["PageName"] = "Team"
		data["Workspace"] = workspace
		team, err := models.GetTeamById(database.DB, uint(team_id))
		if err != nil {
			return ctx.Status(404).SendString("Team not found")
		}
		data["Team"] = team

		available_users, err := models.GetAvailableUsersForTeam(database.DB, team.ID)
		if err != nil {
			return ctx.Status(500).SendString("Error getting available users")
		}

		data["Available_Users"] = available_users
		return ctx.Render("workspaces/team/edit", data, "layouts/workspace")
	})

	app.Get("/:team_id/chat", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)
		workspace_id, _ := ctx.ParamsInt("workspace_id")
		team_id, _ := ctx.ParamsInt("team_id")
		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id))
		team, _ := models.GetTeamById(database.DB, uint(team_id))
		data["Workspace"] = workspace
		data["Team"] = team
		data["PageTitle"] = team.ChatGroup.Name + " Chat"
		data["ChatGroup"] = team.ChatGroup
		data["ChatBoxUsers"] = team.Users

		return ctx.Render("partials/chat", data, "layouts/team")
	})

	app.Post("/:team_id/update", func(ctx *fiber.Ctx) error {
		team_id, _ := ctx.ParamsInt("team_id")
		name := ctx.FormValue("name")

		_, err := models.UpdateTeamName(database.DB, uint(team_id), name)
		if err != nil {
			return ctx.Status(500).SendString("Error updating team")
		}
		return ctx.Render("partials/success-message-with-disappear", fiber.Map{
			"Message": "Team update successfully.",
		})
	})

	app.Post("/:team_id/users/add", func(ctx *fiber.Ctx) error {
		team_id, _ := ctx.ParamsInt("team_id")
		user_id, _ := strconv.Atoi(ctx.FormValue("user_id"))
		workspace_id, _ := ctx.ParamsInt("workspace_id")
		data := load_locals.LoadLocals(ctx)

		err := models.AddUserToTeam(database.DB, uint(team_id), uint(user_id))
		if err != nil {
			log.Println("Error adding user to team: ", err)
			return ctx.Status(500).SendString("Error adding user to team")
		}

		notificaiton_link := fmt.Sprintf("/workspaces/%d/teams/%d", workspace_id, team_id)
		notification := models.Notification{
			Title:  "Added to Team",
			Body:   "You have been added to the team " + strconv.Itoa(team_id),
			UserId: uint(user_id),
			Seen:   false,
			Link:   &notificaiton_link,
		}
		_, err = models.AddNotification(database.DB, &notification)
		if err != nil {
			return ctx.Status(500).SendString("Error creating notification")
		}

		team, _ := models.GetTeamById(database.DB, uint(team_id))
		data["Team"] = team

		workspace, _ := models.GetWorkspaceById(database.DB, uint(workspace_id))
		data["Workspace"] = workspace

		team_user, _ := models.GetUserById(database.DB, uint(user_id))
		data["TeamUser"] = team_user
		return ctx.Render("partials/teams/user", data)
	})

	app.Delete("/:team_id/users/:user_id", func(ctx *fiber.Ctx) error {
		team_id, _ := ctx.ParamsInt("team_id")
		user_id, _ := ctx.ParamsInt("user_id")

		err := models.RemoveUserFromTeam(database.DB, uint(team_id), uint(user_id))
		if err != nil {
			log.Println("Error removing user from team: ", err)
			return ctx.Status(500).SendString("Error removing user from team")
		}

		team, _ := models.GetTeamById(database.DB, uint(team_id))
		notification := models.Notification{
			Title:  "Removed from Team",
			Body:   "You have been removed from the team " + team.Name,
			UserId: uint(user_id),
			Seen:   false,
			Link:   nil,
		}
		_, err = models.AddNotification(database.DB, &notification)
		if err != nil {
			return ctx.Status(500).SendString("Error creating notification")
		}

		return ctx.Render("partials/success-message-with-disappear", fiber.Map{
			"Message": "User removed successfully.",
		})
	})
}
