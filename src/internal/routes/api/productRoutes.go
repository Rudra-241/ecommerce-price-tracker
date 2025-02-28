package api

import (
	"ecommerce-price-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.RouterGroup) {
	r.POST("/product", handlers.CreateProduct)
	r.GET("/products", handlers.GetTrackedProducts)
	r.GET("/product/:id", handlers.GetPriceHistory)
}
