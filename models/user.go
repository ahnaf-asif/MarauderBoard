package models

type User struct {
	ID       int `gorm:"primary_key"`
	Name     string
	Email    string
	Password string
}
