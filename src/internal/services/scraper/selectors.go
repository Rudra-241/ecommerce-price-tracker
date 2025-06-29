package scraper

import (
	"gorm.io/gorm"

	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
)

type Selectors struct {
	Price     string
	Name      string
	Image     string
	ImageAttr string
}

func loadSelectors(store string, defaults Selectors) Selectors {
	var row models.StoreSelector
	if err := db.GetDB().Where("store = ?", store).First(&row).Error; err != nil {
		return defaults
	}
	return Selectors{
		Price:     row.Price,
		Name:      row.Name,
		Image:     row.Image,
		ImageAttr: row.ImageAttr,
	}
}

func SeedSelectors(gdb *gorm.DB) (int, error) {
	created := 0
	for id, s := range registry {
		var count int64
		if err := gdb.Model(&models.StoreSelector{}).Where("store = ?", id).Count(&count).Error; err != nil {
			return created, err
		}
		if count > 0 {
			continue
		}

		d := s.DefaultSelectors()
		row := models.StoreSelector{
			Store:     id,
			Price:     d.Price,
			Name:      d.Name,
			Image:     d.Image,
			ImageAttr: d.ImageAttr,
		}
		if err := gdb.Create(&row).Error; err != nil {
			return created, err
		}
		created++
	}
	return created, nil
}
