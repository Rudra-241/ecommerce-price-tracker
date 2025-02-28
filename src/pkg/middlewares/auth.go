package middlewares

import (
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/register" || c.Request.URL.Path == "/api/login" {
			c.Next()
			return
		}

		tokenString, err := c.Cookie("access-token")
		if err != nil {
			refreshToken, err := c.Cookie("refresh-token")
			if err != nil {
				c.Redirect(http.StatusFound, "/login")
				return
			}
			refreshClaims, err := utils.VerifyToken(refreshToken, utils.RefreshToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
				c.Abort()
				return
			}
			newAccessToken, err := utils.CreateToken(refreshClaims.UserID, refreshClaims.Email, models.Customer, utils.AccessToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new access token"})
				c.Abort()
				return
			}
			c.SetCookie("access-token", newAccessToken, 3600, "/", "localhost", true, true)
			c.Set("userID", refreshClaims.UserID)
			c.Set("email", refreshClaims.Email)
			c.Set("role", refreshClaims.Roles)

			c.Next()
			return
		}

		claims, err := utils.VerifyToken(tokenString, utils.AccessToken)
		if err != nil {
			refreshToken, err := c.Cookie("refresh-token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token and no refresh token available"})
				c.Abort()
				return
			}

			refreshClaims, err := utils.VerifyToken(refreshToken, utils.RefreshToken)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
				c.Abort()
				return
			}
			newAccessToken, err := utils.CreateToken(refreshClaims.UserID, refreshClaims.Email, models.Customer, utils.AccessToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new access token"})
				c.Abort()
				return
			}

			c.SetCookie("access-token", newAccessToken, 3600, "/", "localhost", true, true)

			claims = refreshClaims
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Roles)

		c.Next()
	}
}
