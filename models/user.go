package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name       string
	Email      string
	Password   string
	Workspaces []*Workspace `gorm:"many2many:user_workspaces;"`
	Teams      []*Team      `gorm:"many2many:team_users;"`
}
