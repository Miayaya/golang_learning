package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"size:32;unique;not null"`
	Password  string `gorm:"size:64;not null"`
	Email     string `gorm:"size:64;unique;not null"`
	Posts     []Post
	PostCount int `gorm:"default:0"`
}
