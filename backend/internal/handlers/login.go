package handlers

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/pkg/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	var login models.User
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

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
	var actualUser models.User
	result := tx.Where("email = ?", login.Email).First(&actualUser)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User doesn't exists",
		})
		c.Redirect(http.StatusBadRequest, "/api/register")
		return
	}
	if err := utils.VerifyPassword(actualUser.Password, login.Password); err == nil {
		fmt.Println(login.ID)
		refreshToken, _ := utils.CreateToken(strconv.Itoa(int(actualUser.ID)), actualUser.Email, models.Customer, utils.RefreshToken)
		accessToken, _ := utils.CreateToken(strconv.Itoa(int(actualUser.ID)), actualUser.Email, models.Customer, utils.AccessToken)
		c.SetCookie("refresh-token", refreshToken, 7*24*3600, "/", "localhost", true, true)
		c.SetCookie("access-token", accessToken, 3600, "/", "localhost", true, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
	}
}
