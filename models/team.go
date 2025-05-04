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
	Users       []*User    `gorm:"many2many:team_users;"`
	Projects    []*Project `gorm:"many2many:team_projects;"`
	ChatGroup   *ChatGroup `gorm:"foreignKey:ChatGroupId"`
	ChatGroupId uint
}

func AddNewTeam(db *gorm.DB, team *Team) (*Team, error) {
	if err := db.Create(team).Error; err != nil {
		return nil, err
	}
	return team, nil
}

func GetTeamById(db *gorm.DB, id int) (*Team, error) {
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

func DeleteTeam(db *gorm.DB, team *Team) error {
	if err := db.Delete(team).Error; err != nil {
		return err
	}
	return nil
}
