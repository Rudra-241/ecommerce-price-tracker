package services

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func UpdateAll() {
	dbb := db.GetDB()
	var products []models.Product
	result := dbb.Find(&products)

	if result.Error != nil {
		log.Printf("Error fetching products: %v", result.Error)
		return
	}

	for _, product := range products {
		if err := updateProductInTransaction(dbb, &product); err != nil {
			log.Printf("Failed to update product %d: %v", product.ID, err)
		}
		//added random delays to not get banned
		time.Sleep(time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10)+1) * time.Second)
	}
}

func updateProductInTransaction(db *gorm.DB, product *models.Product) error {
	return db.Transaction(func(tx *gorm.DB) error {
		currentPrice, err := GetPrice(product.Url)
		if err != nil {
			return err
		}

		// Create a new price stamp record
		priceStamp := models.PriceStamp{
			ProductID: product.ID,
			Price:     currentPrice,
			ChangedAt: time.Now(),
		}

		// Save the new price stamp
		if err := tx.Create(&priceStamp).Error; err != nil {
			return err
		}

		// Update product direction based on price changes
		if currentPrice < product.Price {
			product.Direction = models.Decreased
		} else if currentPrice > product.Price {
			product.Direction = models.Increased
		} else {
			product.Direction = models.Unchanged
		}

		//check if prices less than starting price
		var firstPriceStamp models.PriceStamp
		if err := tx.Where("product_id = ?", product.ID).Order("changed_at asc").First(&firstPriceStamp).Error; err == nil {
			if currentPrice < firstPriceStamp.Price {
				product.Direction = models.BelowStart
			}
		}
		product.Price = currentPrice
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		return nil
	})
}
