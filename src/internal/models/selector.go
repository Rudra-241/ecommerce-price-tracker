package models

import "time"

// StoreSelector holds CSS selectors for scraping, keyed by store identifier.
type StoreSelector struct {
	Store     string `gorm:"primaryKey"`
	Price     string `gorm:"not null"`
	Name      string `gorm:"not null"`
	Image     string `gorm:"not null"`
	ImageAttr string `gorm:"not null;default:src"`
	UpdatedAt time.Time
}
