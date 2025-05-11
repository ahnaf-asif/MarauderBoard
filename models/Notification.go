package models

import (
	"log"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model

	Title  string
	Body   string
	UserId uint
	User   User `gorm:"foreignKey:UserId"`
	Seen   bool
	Link   *string
}

func AddNotification(db *gorm.DB, notification *Notification) (*Notification, error) {
	log.Println("Adding notification: ", notification)
	if err := db.Create(notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}

func GetNotificationsByUserId(db *gorm.DB, userId uint) ([]*Notification, error) {
	notifications := []*Notification{}
	if err := db.Where("user_id = ?", userId).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func GetUnseenNotificationsByUserId(db *gorm.DB, userId uint) ([]*Notification, error) {
	notifications := []*Notification{}

	if err := db.Where("user_id = ? AND seen = ?", userId, false).
		Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func MarkNotificationAsSeen(db *gorm.DB, notificationId uint) error {
	if err := db.Model(&Notification{}).
		Where("id = ?", notificationId).Update("seen", true).Error; err != nil {
		return err
	}
	return nil
}
