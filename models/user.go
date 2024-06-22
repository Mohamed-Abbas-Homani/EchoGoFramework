package models

import (
	"gorm.io/gorm"
)

// User Model
type User struct {
	gorm.Model
	Username       string `json:"username" gorm:"unique"`
	Email          string `json:"email" gorm:"unique;not null"`
	Password       string `json:"password" gorm:"not null"`
	ProfilePicture string `json:"profilePicture"`
}
