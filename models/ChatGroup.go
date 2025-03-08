package models

import "gorm.io/gorm"

type ChatGroup struct {
	gorm.Model

	Name     string
	Messages []*ChatMessage `gorm:"foreignKey:ChatGroupId"`
	Users    []*User        `gorm:"many2many:chat_group_users;"`
}
