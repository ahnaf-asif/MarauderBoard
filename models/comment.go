package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	Content  string
	Task     Task `gorm:"foreignKey:TaskId"`
	TaskId   uint
	User     User `gorm:"foreignKey:UserId"`
	UserId   uint
	Replies  []*Comment `gorm:"foreignKey:ParentId"`
	Parent   *Comment   `gorm:"foreignKey:ParentId"`
	ParentId *uint
}

func AddNewComment(db *gorm.DB, comment *Comment) (*Comment, error) {
	if err := db.Create(comment).Error; err != nil {
		return nil, err
	}
	if err := db.Preload("User").First(comment, comment.ID).Error; err != nil {
		return nil, err
	}
	return comment, nil
}
