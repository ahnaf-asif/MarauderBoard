package models

import (
	"log"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	FirstName  string
	LastName   string
	Provider   *string
	Email      string
	Avatar     *string
	Workspaces []*Workspace `gorm:"many2many:user_workspaces;"`
	Teams      []*Team      `gorm:"many2many:team_users;"`
}

func AddNewUser(db *gorm.DB, user User) (*User, error) {
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	log.Println("Getting user by email:", email)
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *User) error {
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
