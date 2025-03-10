package handlers

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/pkg/utils"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var registration models.User

	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	//ISSUE: doesn't check if email is valid or not
	if !isValidEmail(registration.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address",
		})
	}
	if !isValidPassword(registration.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
	}
	hashedPassword, err := utils.HashPassword(registration.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	registration.Password = hashedPassword

	dbb := db.GetDB()

	tx := dbb.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start database transaction",
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingUser models.User
	result := tx.Where("email = ?", registration.Email).First(&existingUser)
	if result.Error == nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
		})
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}
	registration.Role = models.Customer
	if err := tx.Create(&registration).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}
	fmt.Print(registration.ID)

	refreshToken, _ := utils.CreateToken(strconv.Itoa(int(registration.ID)), registration.Email, models.Customer, utils.RefreshToken)
	accessToken, _ := utils.CreateToken(strconv.Itoa(int(registration.ID)), registration.Email, models.Customer, utils.AccessToken)
	c.SetCookie("refresh-token", refreshToken, 7*24*3600, "/", "localhost", true, true)
	c.SetCookie("access-token", accessToken, 3600, "/", "localhost", true, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	return len(password) >= 3
}
