package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie("access-token", "", -1, "/", "localhost", true, true)
	c.SetCookie("refresh-token", "", -1, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func Me(c *gin.Context) {
	id, _ := c.Get("userID")
	email, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{"id": id, "email": email})
}
