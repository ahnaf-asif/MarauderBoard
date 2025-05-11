package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	FirstName  string
	LastName   string
	Provider   *string
	Email      string `gorm:"unique;not null"`
	Password   *string
	Avatar     *string
	Workspaces []*Workspace `gorm:"many2many:user_workspaces;"`
	Teams      []*Team      `gorm:"many2many:team_users;"`
}

func GeneratePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

func GetUserById(db *gorm.DB, id uint) (*User, error) {
	var user User
	if err := db.Preload("Workspaces").Preload("Teams").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers(db *gorm.DB) ([]*User, error) {
	var users []*User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func UpdateUser(db *gorm.DB, user *User) (*User, error) {
	if err := db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
