package handlers

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/pkg/utils"
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

	db := db.GetDB()

	tx := db.Begin()
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
		refresh_token, _ := utils.CreateToken(strconv.Itoa(int(login.ID)), login.Email, models.Customer, utils.RefreshToken)
		access_token, _ := utils.CreateToken(strconv.Itoa(int(login.ID)), login.Email, models.Customer, utils.AccessToken)
		c.SetCookie("refresh-token", refresh_token, 7*24*3600, "/", "localhost", true, true)
		c.SetCookie("access-token", access_token, 3600, "/", "localhost", true, true)
		c.Redirect(http.StatusOK, "/api/testing")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
	}
}
