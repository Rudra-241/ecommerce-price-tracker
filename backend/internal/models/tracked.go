package models

import (
	"time"
)

type Tracked struct {
	UserID    uint      `gorm:"primaryKey"`
	ProductID uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
