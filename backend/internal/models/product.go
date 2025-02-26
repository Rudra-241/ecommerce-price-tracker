package models

import (
	"gorm.io/gorm"
	"time"
)

//TODO: add `json:"fieldname"` thing to every attribute

type PriceState string

const (
	Increased  = "Increased"
	Decreased  = "Decreased"
	Unchanged  = "Unchanged"
	BelowStart = "BelowStart"
)

func (p *PriceState) Scan(value interface{}) error {
	*p = PriceState(value.(string))
	return nil
}

func (p *Product) AfterUpdate(tx *gorm.DB) (err error) {
	p.ModifiedAt = time.Now()
	return nil
}

type PriceStamp struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `gorm:"not null;index"`
	Price     float64   `gorm:"not null"`
	ChangedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type Product struct {
	ID            uint    `gorm:"primaryKey"` // TODO: change primary key
	Name          string  `gorm:"not null"`
	Price         float64 `gorm:"not null"`
	ImgLink       string
	Url           string
	PriceHistory  []PriceStamp `gorm:"foreignKey:ProductID"`
	UsersTracking []User       `gorm:"many2many:tracked;"`
	Direction     PriceState   `gorm:"default:Unchanged;not null"`
	ModifiedAt    time.Time    `gorm:"default:CURRENT_TIMESTAMP"`
}

type ProductInfo struct {
	Name    string
	Price   float64
	ImgLink string
	Url     string
}

type ProductRequest struct {
	Url string `json:"url"`
}
