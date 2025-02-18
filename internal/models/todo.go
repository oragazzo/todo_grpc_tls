package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents a todo item in the database
type Todo struct {
	gorm.Model
	Title     string    `gorm:"not null"`
	Completed bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
}
