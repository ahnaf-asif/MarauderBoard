package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	Name        string
	Workspace   Workspace `gorm:"foreignKey:WorkspaceId"`
	WorkspaceId int
	Leader      User `gorm:"foreignKey:LeaderId"`
	LeaderId    int
	Users       []*User    `gorm:"many2many:team_users;"`
	Projects    []*Project `gorm:"many2many:team_projects;"`
	ChatGroup   *ChatGroup `gorm:"foreignKey:ChatGroupId"`
	ChatGroupId int
}
