package services

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/internal/services/scraper"
	"log/slog"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

const (
	recentPriceWindow   = 30 * 24 * time.Hour
	minRecentSamples    = 5
	significantDropFrac = 0.05
)

func UpdateAll() {
	dbb := db.GetDB()
	var products []models.Product
	result := dbb.Find(&products)

	if result.Error != nil {
		slog.Error("fetching products", "err", result.Error)
		return
	}

	for _, product := range products {
		if err := updateProductInTransaction(dbb, &product); err != nil {
			slog.Error("failed to update product", "product_id", product.ID, "err", err)
		}
		//added random delays to not get banned
		time.Sleep(time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10)+1) * time.Second)
	}
}

func updateProductInTransaction(db *gorm.DB, product *models.Product) error {
	return db.Transaction(func(tx *gorm.DB) error {
		s, err := scraper.For(product.Url)
		if err != nil {
			return err
		}
		currentPrice, err := s.GetPrice(product.Url)
		if err != nil {
			return err
		}

		var recent struct {
			Avg   float64
			Count int64
		}
		if err := tx.Model(&models.PriceStamp{}).
			Where("product_id = ? AND changed_at >= ?", product.ID, time.Now().Add(-recentPriceWindow)).
			Select("AVG(price) AS avg, COUNT(*) AS count").
			Scan(&recent).Error; err != nil {
			return err
		}

		priceStamp := models.PriceStamp{
			ProductID: product.ID,
			Price:     currentPrice,
			ChangedAt: time.Now(),
		}

		if err := tx.Create(&priceStamp).Error; err != nil {
			return err
		}

		if currentPrice < product.Price {
			product.Direction = models.Decreased
		} else if currentPrice > product.Price {
			product.Direction = models.Increased
		} else {
			product.Direction = models.Unchanged
		}

		// var firstPriceStamp models.PriceStamp
		// if err := tx.Where("product_id = ?", product.ID).Order("changed_at asc").First(&firstPriceStamp).Error; err == nil {
		// 	if currentPrice < firstPriceStamp.Price {
		// 		product.Direction = models.BelowStart
		// 	}
		// }

		if recent.Count >= minRecentSamples && currentPrice <= recent.Avg*(1-significantDropFrac) {
			product.Direction = models.BelowRecentAvg
		}
		product.Price = currentPrice
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		return nil
	})
}
