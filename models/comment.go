package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	Content  string
	Task     Task `gorm:"foreignKey:TaskId"`
	TaskId   int
	User     User `gorm:"foreignKey:UserId"`
	UserId   int
	Replies  []*Comment `gorm:"foreignKey:ParentId"`
	Parent   *Comment   `gorm:"foreignKey:ParentId"`
	ParentId *int
}
