package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	Name        string
	Description string
	Status      string
	Workspace   Workspace `gorm:"foreignKey:WorkspaceId"`
	WorkspaceId int
	Teams       []*Team `gorm:"many2many:team_projects;"`
}
