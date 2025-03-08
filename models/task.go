package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model

	Name        string
	Description string
	Status      string
	Project     Project `gorm:"foreignKey:ProjectId"`
	ProjectId   int
	Assignee    User `gorm:"foreignKey:AssigneeId"`
	AssigneeId  int
	Reporter    User `gorm:"foreignKey:ReporterId"`
	ReporterId  int
	Team        Team `gorm:"foreignKey:TeamId"`
	TeamId      int
	Comments    []*Comment `gorm:"foreignKey:TaskId"`
}
