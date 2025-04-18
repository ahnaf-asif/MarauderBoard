package models

import "gorm.io/gorm"

type ChatGroup struct {
	gorm.Model

	Name     string
	Messages []*ChatMessage `gorm:"foreignKey:ChatGroupId"`
	Users    []*User        `gorm:"many2many:chat_group_users;"`
}

func AddChatGroup(db *gorm.DB, chatGroup *ChatGroup) (*ChatGroup, error) {
	if err := db.Create(chatGroup).Error; err != nil {
		return nil, err
	}
	return chatGroup, nil
}
