package model

import (
	"gorm.io/gorm"
)

// type User struct {
// 	gorm.Model
// 	UserName  string `gorm:"uniqueIndex"`
// 	Email     string `gorm:"uniqueIndex"`
// 	Password  string
// 	FirstName string
// 	LastName  string
// 	IsAdmin   bool
// }

type Todo struct {
	gorm.Model
	OwnerID uint
	Title   string
	Content string
	Done    bool
}
