package models

import (
	"gorm.io/gorm"
)

// User Model
type User struct {
	gorm.Model
	Email          string `json:"email" gorm:"unique;not null"`
	Password       string `json:"password" gorm:"not null"`
	ProfilePicture string `json:"profilePicture"`
}
