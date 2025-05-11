package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	Name        string
	Workspace   Workspace `gorm:"foreignKey:WorkspaceId"`
	WorkspaceId uint
	Leader      User `gorm:"foreignKey:LeaderId"`
	LeaderId    uint
	Users       []*User    `gorm:"many2many:team_users;constraint:OnDelete:SET NULL;"`
	Projects    []*Project `gorm:"many2many:team_projects;constraint:OnDelete:SET NULL;"`
	ChatGroup   *ChatGroup `gorm:"foreignKey:ChatGroupId;constraint:OnDelete:SET NULL;"`
	ChatGroupId uint
}

func AddNewTeam(db *gorm.DB, team *Team) (*Team, error) {
	if err := db.Create(team).Error; err != nil {
		return nil, err
	}
	return team, nil
}

func GetTeamById(db *gorm.DB, id uint) (*Team, error) {
	team := &Team{}
	if err := db.Preload("Leader").
		Preload("Users").
		Preload("Projects").
		Preload("ChatGroup").
		First(team, id).Error; err != nil {
		return nil, err
	}
	return team, nil
}

func GetTeamsByWorkspaceId(db *gorm.DB, workspaceId uint) ([]*Team, error) {
	teams := []*Team{}
	if err := db.Preload("Leader").
		Preload("Users").
		Preload("Projects").
		Preload("ChatGroup").
		Preload("Workspace").
		Where("workspace_id = ?", workspaceId).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func GetTeamsByWorkspaceAndUserId(db *gorm.DB, workspaceId uint, userId uint) ([]*Team, error) {
	teams := []*Team{}
	if err := db.Preload("Leader").
		Preload("Users").
		Preload("Projects").
		Preload("ChatGroup").
		Joins("JOIN team_users ON teams.id = team_users.team_id").
		Where("team_users.user_id = ?", userId).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func GetTeamsByUserId(db *gorm.DB, user_id int) ([]*Team, error) {
	teams := []*Team{}
	if err := db.Preload("Leader").
		Preload("Users").
		Preload("Projects").
		Preload("ChatGroup").
		Joins("JOIN team_users ON teams.id = team_users.team_id").
		Where("team_users.user_id = ?", user_id).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func UpdateTeam(db *gorm.DB, team *Team) (*Team, error) {
	if err := db.Save(team).Error; err != nil {
		return nil, err
	}
	return team, nil
}

func UpdateTeamName(db *gorm.DB, teamId uint, name string) (*Team, error) {
	if err := db.Model(&Team{}).Where("id = ?", teamId).Update("name", name).Error; err != nil {
		return nil, err
	}
	return GetTeamById(db, teamId)
}

func DeleteTeam(db *gorm.DB, team *Team) error {
	if err := db.Delete(team).Error; err != nil {
		return err
	}
	return nil
}

func AddUserToTeam(db *gorm.DB, tream_id uint, user_id uint) error {
	if err := db.Exec("INSERT INTO team_users (team_id, user_id) VALUES (?, ?)", tream_id, user_id).Error; err != nil {
		return err
	}
	return nil
}

func RemoveUserFromTeam(db *gorm.DB, team_id uint, user_id uint) error {
	if err := db.Exec("DELETE FROM team_users WHERE team_id = ? AND user_id = ?", team_id, user_id).Error; err != nil {
		return err
	}
	return nil
}

func GetAvailableUsersForTeam(db *gorm.DB, teamId uint) ([]*User, error) {
	all_users, _ := GetAllUsers(db)
	team, _ := GetTeamById(db, teamId)

	available_users := []*User{}
	for _, user := range all_users {
		take := true
		for _, team_user := range team.Users {
			if user.ID == team_user.ID {
				take = false
				break
			}
		}

		if take {
			available_users = append(available_users, user)
		}
	}
	return available_users, nil
}

func GetAvailableTeamsForProject(db *gorm.DB, workspaceId uint, projectId uint) ([]*Team, error) {
	teams := []*Team{}
	all_teams, _ := GetTeamsByWorkspaceId(db, workspaceId)
	project, _ := GetProjectById(db, projectId)
	for _, team := range all_teams {
		take := true
		for _, project_team := range project.Teams {
			if team.ID == project_team.ID {
				take = false
				break
			}
		}
		if take {
			teams = append(teams, team)
		}
	}

	return teams, nil
}
