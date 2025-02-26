package routes

import (
	"ecommerce-price-tracker/internal/handlers"
	middleware "ecommerce-price-tracker/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	api.POST("/register", handlers.CreateUser)
	api.POST("/login", handlers.LoginUser)
	api.POST("/product", handlers.CreateProduct)
	api.GET("/products", handlers.GetTrackedProducts)
	api.GET("/product/:id", handlers.GetPriceHistory)
}
