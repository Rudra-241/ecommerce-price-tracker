package models

import (
	"time"
)

type Role string

const (
	Customer Role = "user"
	Admin    Role = "admin"
)

type User struct {
	ID              uint      `gorm:"primaryKey"` // TODO: change primary key
	Email           string    `gorm:"uniqueIndex;not null"`
	Password        string    `gorm:"not null"`
	Role            Role      `gorm:"not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
	TrackedProducts []Product `gorm:"many2many:tracked;"`
}
