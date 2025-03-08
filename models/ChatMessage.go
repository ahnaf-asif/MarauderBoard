package models

import "gorm.io/gorm"

type ChatMessage struct {
	gorm.Model

	Content     string
	User        User `gorm:"foreignKey:UserId"`
	UserId      int
	ChatGroup   ChatGroup `gorm:"foreignKey:ChatGroupId"`
	ChatGroupId int
}
