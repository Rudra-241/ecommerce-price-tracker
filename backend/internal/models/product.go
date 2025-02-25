package models

type Product struct {
	ID            uint    `gorm:"primaryKey"`
	Name          string  `gorm:"not null"`
	Price         float64 `gorm:"not null"`
	ImgLink       string
	Url           string
	UsersTracking []User `gorm:"many2many:tracked;"`
}
