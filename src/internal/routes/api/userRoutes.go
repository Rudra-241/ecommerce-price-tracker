package api

import (
	"ecommerce-price-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.POST("/register", handlers.CreateUser)
	r.POST("/login", handlers.LoginUser)
	r.POST("/logout", handlers.Logout)
	r.GET("/me", handlers.Me)
}
