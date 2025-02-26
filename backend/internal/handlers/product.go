package handlers

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateProduct(c *gin.Context) {
	var product models.Product

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
		panic("User not found in context")
		return
	}

	dbb := db.GetDB()

	tx := dbb.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
		})
	}

	var user models.User

	result := tx.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User doesn't exists",
		})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

}
