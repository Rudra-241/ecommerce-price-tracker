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

func GetPriceHistory(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product ID is required",
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	dbb := db.GetDB()

	var user models.User
	if err := dbb.Preload("TrackedProducts", "id = ?", productID).Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	//check if user is currently tracking the product or not

	if len(user.TrackedProducts) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "User is not tracking this product",
		})
		return
	}

	var product models.Product
	if err := dbb.Preload("PriceHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("changed_at ASC")
	}).First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch product data",
			})
		}
		return
	}

	response := gin.H{
		"product_id":    product.ID,
		"product_name":  product.Name,
		"current_price": product.Price,
		"price_history": product.PriceHistory,
	}

	c.JSON(http.StatusOK, response)
}
