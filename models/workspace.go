package models

import "gorm.io/gorm"

type Workspace struct {
	gorm.Model

	Name            string
	Description     string
	Administrator   *User `gorm:"foreignKey:AdministratorId"`
	AdministratorId *int
	Users           []*User    `gorm:"many2many:user_workspaces;"`
	Teams           []*Team    `gorm:"foreignKey:WorkspaceId"`
	Projects        []*Project `gorm:"foreignKey:WorkspaceId"`
	ChatGroup       *ChatGroup `gorm:"foreignKey:ChatGroupId"`
	ChatGroupId     *int
}
