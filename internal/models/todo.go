package models

import (
	"gorm.io/gorm"
)

// Todo represents a todo item in the database
type Todo struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Completed bool   `gorm:"not null"`
}
