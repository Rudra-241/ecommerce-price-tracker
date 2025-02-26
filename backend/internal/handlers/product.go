package handlers

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func CreateProduct(c *gin.Context) {
	var product models.ProductRequest

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found in context",
		})
		return
	}

	url := product.Url
	productInfo, err := services.GetProductInfo(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get product info",
		})
		return
	}

	dbb := db.GetDB()

	tx := dbb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
		return
	}

	var user models.User
	result := tx.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User doesn't exist",
		})
		return
	}

	var existingProduct models.Product
	result = tx.Where("url = ?", productInfo.Url).First(&existingProduct)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newProduct := models.Product{
			Url:          productInfo.Url,
			ImgLink:      productInfo.ImgLink,
			Name:         productInfo.Name,
			Price:        productInfo.Price,
			PriceHistory: make([]models.PriceStamp, 0),
		}

		if err := tx.Create(&newProduct).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create product",
			})
			return
		}

		priceStamp := models.PriceStamp{
			ProductID: newProduct.ID,
			Price:     productInfo.Price,
			ChangedAt: time.Now(),
		}

		if err := tx.Create(&priceStamp).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create price stamp",
			})
			return
		}

		existingProduct = newProduct
	} else if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error when checking for existing product",
		})
		return
	}
	var count int64
	if err := tx.Table("tracked").
		Where("user_id = ? AND product_id = ?", user.ID, existingProduct.ID).
		Count(&count).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error while checking tracked products",
		})
		return
	}
	if count > 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You are already tracking this product",
		})
		return
	}

	if err := tx.Model(&user).Association("TrackedProducts").Append(&existingProduct); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to associate product with user",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product tracked successfully",
	})
}

func GetTrackedProducts(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found in context",
		})
		return
	}

	dbb := db.GetDB()

	var user models.User
	if err := dbb.Preload("TrackedProducts").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tracked_products": user.TrackedProducts,
	})
}
